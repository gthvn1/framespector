package network

import (
	"fmt"
	"log/slog"
	"net"
	"os/exec"

	"golang.org/x/sys/unix"
)

type Veth struct {
	P1name string
	P2name string
	IP     string
	FD     int
	Logger *slog.Logger
}

// htons() function converts the unsigned short integer "hostshort"
// from host byte order to network byte order.
func htons(i uint16) uint16 {
	return (i<<8 | i>>8)
}

func NewVeth(logger *slog.Logger, name string, ip string) *Veth {
	return &Veth{
		P1name: name,
		P2name: name + "-peer",
		IP:     ip,
		FD:     -1,
		Logger: logger,
	}
}

// On Linux: man veth
func (v *Veth) Setup() error {
	cmd := exec.Command("ip", "link", "add", v.P1name, "type", "veth", "peer", "name", v.P2name)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run command %q, error: %w, output: %s", cmd.String(), err, output)
	}

	// Just return in case of error when setting links up.

	if exec.Command("ip", "link", "set", v.P1name, "up").Run() != nil {
		v.Cleanup()
		return fmt.Errorf("failed to set link %s up", v.P1name)
	}

	if exec.Command("ip", "link", "set", v.P2name, "up").Run() != nil {
		v.Cleanup()
		return fmt.Errorf("failed to set link %s up", v.P2name)
	}

	if exec.Command("ip", "addr", "add", v.IP, "dev", v.P1name).Run() != nil {
		v.Cleanup()
		return fmt.Errorf("failed to add %s to %s", v.IP, v.P1name)
	}

	return nil
}

func (v *Veth) Cleanup() {
	// Just report failure and continue
	if exec.Command("ip", "link", "set", v.P1name, "down").Run() != nil {
		v.Logger.Error("failed to set link down", "veth", v.P1name)
	}

	if exec.Command("ip", "link", "del", v.P1name).Run() != nil {
		v.Logger.Error("failed to delete link", "veth", v.P1name)
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
	return nil
}

func (v *Veth) BindPeer() error {
	if v.FD < 0 {
		return fmt.Errorf("socket is not created")
	}

	iface, err := net.InterfaceByName(v.P2name)
	if err != nil {
		return fmt.Errorf("failed to get peer interface: %w", err)
	}

	// man sockaddr
	sll := &unix.SockaddrLinklayer{
		Protocol: htons(unix.ETH_P_ALL),
		Ifindex:  iface.Index,
	}

	if err := unix.Bind(v.FD, sll); err != nil {
		return fmt.Errorf("failed to bind socket: %w", err)
	}

	return nil
}
