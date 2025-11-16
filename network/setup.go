package network

import (
	"fmt"
	"log/slog"
	"net"
	"os/exec"

	"golang.org/x/sys/unix"
)

type Veth struct {
	HostName string
	PeerName string
	HostIP   net.IP
	HostNet  *net.IPNet
	PeerIP   net.IP
	PeerNet  *net.IPNet
	FD       int
	SAddr    *unix.SockaddrLinklayer
	Logger   *slog.Logger
}

// htons() function converts the unsigned short integer "hostshort"
// from host byte order to network byte order.
func htons(i uint16) uint16 {
	return (i<<8 | i>>8)
}

func stringToIPv4(ipStr string) (net.IP, *net.IPNet, error) {
	ip, ipNet, err := net.ParseCIDR(ipStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid IPv4 %s: %w", ipStr, err)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil, nil, fmt.Errorf("invalid IPv4 %s", ipStr)
	}

	return ipv4, ipNet, nil
}

func NewVeth(logger *slog.Logger, name string, host string, peer string) (*Veth, error) {

	HostIP, HostNet, err1 := stringToIPv4(host)
	if err1 != nil {
		return nil, err1
	}

	PeerIP, PeerNet, err2 := stringToIPv4(peer)
	if err2 != nil {
		return nil, err2
	}

	return &Veth{
		HostName: name,
		PeerName: name + "-peer",
		HostIP:   HostIP,
		HostNet:  HostNet,
		PeerIP:   PeerIP,
		PeerNet:  PeerNet,
		FD:       -1,
		SAddr:    nil,
		Logger:   logger,
	}, nil
}

// On Linux: man veth
func (v *Veth) Setup() error {
	cmd := exec.Command("ip", "link", "add", v.HostName, "type", "veth", "peer", "name", v.PeerName)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run command %q, error: %w, output: %s", cmd.String(), err, output)
	}

	// Just return in case of error when setting links up.

	if exec.Command("ip", "link", "set", v.HostName, "up").Run() != nil {
		v.Cleanup()
		return fmt.Errorf("failed to set link %s up", v.HostName)
	}

	if exec.Command("ip", "link", "set", v.PeerName, "up").Run() != nil {
		v.Cleanup()
		return fmt.Errorf("failed to set link %s up", v.PeerName)
	}

	if exec.Command("ip", "addr", "add", v.HostNet.String(), "dev", v.HostName).Run() != nil {
		v.Cleanup()
		return fmt.Errorf("failed to add %s to %s", v.HostNet.String(), v.HostName)
	}

	return nil
}

func (v *Veth) Cleanup() {
	// Just report failure and continue
	if exec.Command("ip", "link", "set", v.HostName, "down").Run() != nil {
		v.Logger.Error("failed to set link down", "veth", v.HostName)
	}

	if exec.Command("ip", "link", "del", v.HostName).Run() != nil {
		v.Logger.Error("failed to delete link", "veth", v.HostName)
	}

	if v.FD >= 0 {
		if err := unix.Close(v.FD); err != nil {
			v.Logger.Error("failed to close the socket")
		}
		v.FD = -1
	}
}

func (v *Veth) CreateSocket() error {
	// On Linux: man packet
	// => Set protocol to ETH_P_ALL to receive all protocols
	proto := htons(unix.ETH_P_ALL)
	v.Logger.Debug("proto set", "proto", proto)
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(proto))
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}

	v.FD = fd
	v.Logger.Debug("virtual pair socket created")
	return nil
}

func (v *Veth) BindPeer() error {
	if v.FD < 0 {
		return fmt.Errorf("socket is not created")
	}

	iface, err := net.InterfaceByName(v.PeerName)
	if err != nil {
		return fmt.Errorf("failed to get peer interface: %w", err)
	}

	// man sockaddr
	sll := &unix.SockaddrLinklayer{
		Protocol: htons(unix.ETH_P_ALL),
		Ifindex:  iface.Index,
	}

	v.SAddr = sll

	if err := unix.Bind(v.FD, sll); err != nil {
		return fmt.Errorf("failed to bind socket: %w", err)
	}

	v.Logger.Debug("bind done", "iface", v.PeerName)
	return nil
}
