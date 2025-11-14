package network

import (
	"fmt"
	"log"
	"os/exec"
)

type Veth struct {
	P1name string
	P2name string
}

func NewVeth(name string) *Veth {
	return &Veth{
		P1name: name,
		P2name: name + "-peer",
	}
}

// man veth
func (v *Veth) Setup() error {
	cmd := exec.Command("ip", "link", "add", v.P1name, "type", "veth", "peer", "name", v.P2name)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to run command %q\nerror: %w\noutput: %s", cmd.String(), err, output)
	}

	// Just return in case of error when setting links up.

	if exec.Command("ip", "link", "set", v.P1name, "up").Run() != nil {
		v.Cleanup()
		return fmt.Errorf("Failed to set link %s up\n", v.P1name)
	}

	if exec.Command("ip", "link", "set", v.P2name, "up").Run() != nil {
		v.Cleanup()
		return fmt.Errorf("Failed to set link %s up\n", v.P2name)
	}

	return nil
}

func (v *Veth) Cleanup() {
	// Only need to cleanup one points
	if exec.Command("ip", "link", "set", v.P1name, "down").Run() != nil {
		log.Printf("Failed to set link %s down\n", v.P1name)
	}

	if exec.Command("ip", "link", "del", v.P1name).Run() != nil {
		log.Printf("Failed to delete link %s\n", v.P1name)
	}
}
