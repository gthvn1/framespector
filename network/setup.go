package network

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"golang.org/x/sys/unix"
)

type Veth struct {
	P1name string
	P2name string
	FD     int
}

// htons() function converts the unsigned short integer "hostshort"
// from host byte order to network byte order.
func htons(i uint16) uint16 {
	return (i<<8 | i>>8)
}

func NewVeth(name string) *Veth {
	return &Veth{
		P1name: name,
		P2name: name + "-peer",
		FD:     -1,
	}
}

// On Linux: man veth
func (v *Veth) Setup() error {
	cmd := exec.Command("ip", "link", "add", v.P1name, "type", "veth", "peer", "name", v.P2name)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run command %q\nerror: %w\noutput: %s", cmd.String(), err, output)
	}

	// Just return in case of error when setting links up.

	if exec.Command("ip", "link", "set", v.P1name, "up").Run() != nil {
		v.Cleanup()
		return fmt.Errorf("failed to set link %s up\n", v.P1name)
	}

	if exec.Command("ip", "link", "set", v.P2name, "up").Run() != nil {
		v.Cleanup()
		return fmt.Errorf("failed to set link %s up\n", v.P2name)
	}

	return nil
}

func (v *Veth) Cleanup() {
	// Just report failure and continue
	if exec.Command("ip", "link", "set", v.P1name, "down").Run() != nil {
		log.Printf("failed to set link %s down\n", v.P1name)
	}

	if exec.Command("ip", "link", "del", v.P1name).Run() != nil {
		log.Printf("failed to delete link %s\n", v.P1name)
	}

	if v.FD >= 0 {
		if err := unix.Close(v.FD); err != nil {
			log.Printf("failed to close the socket")
		}
		v.FD = -1
	}
}

func (v *Veth) CreateSocket() error {
	// On Linux: man packet
	// => Set protocol to ETH_P_ALL to receive all protocols
	proto := htons(unix.ETH_P_ALL)
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
