package network

import (
	"fmt"
	"log"
	"os/exec"
)

func SetupNetwork(name string) error {
	peer := name + "-peer"

	cmd := exec.Command("ip", "link", "add", name, "type", "veth", "peer", "name", peer)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to run command %q\nerror: %w\noutput: %s", cmd.String(), err, output)
	}

	// Just return in case of error when setting links up.

	if exec.Command("ip", "link", "set", name, "up").Run() != nil {
		return fmt.Errorf("Failed to set link %s down\n", name)
	}

	if exec.Command("ip", "link", "set", peer, "up").Run() != nil {
		return fmt.Errorf("Failed to set link %s down\n", peer)
	}

	return nil
}

func CleanupNetorwk(name string) {
	if exec.Command("ip", "link", "set", name, "down").Run() != nil {
		log.Printf("Failed to set link %s down\n", name)
	}

	if exec.Command("ip", "link", "del", name).Run() != nil {
		log.Printf("Failed to delete link %s\n", name)
	}
}
