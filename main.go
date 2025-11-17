package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"example.com/framespector/network"
	"golang.org/x/sys/unix"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))

	args := ReadArgs()
	if args == nil {
		return
	}

	vethConf := network.VethConf{
		Name:      args.vethName,
		HostIPStr: args.hostIPStr,
		PeerIPStr: args.peerIPStr,
	}

	veth, err := network.NewVeth(logger, vethConf)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := veth.Setup(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer veth.Cleanup()

	if err := veth.CreateSocket(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := veth.BindPeer(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// At this point all fields of Veth are initialized
	if veth.SAddr == nil {
		panic("At this point SAddr should be initialized")
	}

	logger.Info("Setup network done")

	// To be able to quit the loop using ctrl-c we create a channel
	// of type os.Signal with a size of 1
	sigChan := make(chan os.Signal, 1)
	// If ctrl-c is hit sends it to the channel sigChan
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("Hit ctrl-c to quit")

	// start a go routine that will listen on socket
	ctx, cancel := context.WithCancel(context.Background())

	// We need to wait for the go routine to end before closing
	// socket. So we use WaitGroup to track the go routine
	var wg sync.WaitGroup
	wg.Add(1)
	go receiveLoop(ctx, &wg, veth)

	// and block until ctrl-c is received
	<-sigChan
	logger.Info("ctrl-c received, shutting down...")

	// Cancel go routine
	cancel()

	wg.Wait()
	logger.Info("clean shutdown complete")
}

func receiveLoop(ctx context.Context, wg *sync.WaitGroup, veth *network.Veth) {
	// When done signal it
	defer wg.Done()

	// We need to poll to avoid blocking on Recvfrom
	pollFds := []unix.PollFd{
		{
			Fd:     int32(veth.FD),
			Events: unix.POLLIN,
		},
	}

	rawFrame := make([]byte, 4096)

	for {
		select {
		case <-ctx.Done():
			veth.Logger.Info("stop receiving frame")
			return
		default:
			// Poll with a timeout of 100ms
			n, err := unix.Poll(pollFds, 100)
			if err != nil {
				veth.Logger.Warn("poll error", "err", err)
				continue
			}

			if n == 0 {
				// We hit the timeout so just continue and we will be able to
				// check ctx.Done
				continue
			}

			n, _, err = unix.Recvfrom(veth.FD, rawFrame, 0)
			if err == unix.EBADF || err == unix.EINVAL {
				veth.Logger.Error("socket closed")
				return
			}

			if err != nil {
				veth.Logger.Warn("receive error", "err", err)
				continue
			}

			veth.Logger.Info("frame received", "bytes", n)
			f, err := network.ParseEthernet(rawFrame[:n])
			if err != nil {
				veth.Logger.Error("failed to decode frame", "err", err)
				continue
			}

			veth.Logger.Debug(f.String())

			// Dispatch based on the ethernet type
			switch f.EtherType {
			case network.EtherTypeARP:
				handleARP(veth, f.Payload)
			case network.EtherTypeIPv4:
				handleIPv4(veth, f.Payload)
			case network.EtherTypeIPv6:
				handleIPv6(veth, f.Payload)
			case network.EtherTypeVLAN:
				veth.Logger.Debug("VLAN frame ignored")
			case network.EtherTypeUnknown:
				veth.Logger.Warn("unkown ether type", "type", fmt.Sprintf("0x%04x", f.EtherType))
			default:
				// If you are here it is because you modified the EtherType enum and you
				// don't handle it here.
				panic(fmt.Sprintf("unhandled EtherType in switch: 0x%04x", f.EtherType))
			}
		}
	}
}

// ------------------------------------------------------------------------------
// ETHERNET FRAME HANDLERS
func handleIPv6(veth *network.Veth, payload []byte) {
	_, _ = veth, payload
	veth.Logger.Debug("TODO: decode ipv6")

}

func handleIPv4(veth *network.Veth, payload []byte) {
	p, err := network.ParseIPv4Packet(payload, veth.PeerIP)
	if err != nil {
		veth.Logger.Error("failed to parse IPv4 packet", "err", err)
		return
	}

	switch p.Protocol {
	case network.ICMPProtocol:
		icmp, err := network.ParseICMP(p)
		if err != nil {
			veth.Logger.Error("failed to parse ICMP packet", "err", err)
			return
		}

		if icmp.Type != network.ICMPEchoRequest {
			veth.Logger.Warn("only ICMP Echo request are handled")
			return
		}

		veth.Logger.Warn("TODO: handle ICMP echo request")
	default:
		veth.Logger.Warn("Only ICMP protocol is managed currently")
	}
}

func handleARP(veth *network.Veth, payload []byte) {
	peerIface, err1 := net.InterfaceByName(veth.PeerName)
	if err1 != nil {
		veth.Logger.Error("failed to peer interface", "iface", veth.PeerName, "err", err1)
		return
	}

	reply, err2 := network.ParseARP(payload, peerIface.HardwareAddr, veth.PeerIP)
	if err2 != nil {
		veth.Logger.Error("ARP request not handled", "err", err2)
		return
	}

	arpPayload := reply.Marshal()
	ethFrame := network.BuildEthernetFrame(reply.TargetHA, reply.SenderHA, network.EtherTypeARP, arpPayload)

	// TODO: better handling of veth.SAddr, currently we don't check if it is nil, but we
	//       maybe need to add a method to send it in ethernet.go
	if err := unix.Sendto(veth.FD, ethFrame, 0, veth.SAddr); err != nil {
		veth.Logger.Error("failed to send ARP reply", "err", err)
	} else {
		veth.Logger.Info("sent ARP reply", "to_mac", reply.TargetHA.String(), "from_mac", reply.SenderHA.String())
	}
}

// ------------------------------------------------------------------------------
// READ ARGUMENTS
type Args struct {
	vethName  string
	hostIPStr string
	peerIPStr string
}

func ReadArgs() *Args {
	// We are expecting --veth <ifacename> and --ip <x.x.x.x/yy>
	// So the virtual pair name and the ip with its subnet
	vethName := flag.String("veth", "veth0", "Virtual Pair name")
	hostIP := flag.String("ip", "192.168.35.2/24", "IP address with CIDR")
	peerIP := flag.String("peer", "192.168.35.3/24", "IP address of the peer with CIDR")
	help := flag.Bool("help", false, "Print help")

	flag.Parse()

	if *help {
		fmt.Println("Usage: framespector --veth <veth-name> --ip <ip/cidr> --peer <ip/cidr>")
		flag.PrintDefaults()
		return nil
	}

	// Just check that IPs are valid
	if _, _, err := net.ParseCIDR(*hostIP); err != nil {
		fmt.Printf("%s is not a valid IP address with CIDR\n", *hostIP)
		return nil
	}

	if _, _, err := net.ParseCIDR(*peerIP); err != nil {
		fmt.Printf("%s is not a valid IP address with CIDR\n", *peerIP)
		return nil
	}

	return &Args{
		vethName:  *vethName,
		hostIPStr: *hostIP,
		peerIPStr: *peerIP,
	}
}
