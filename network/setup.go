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

	log.Printf("Successfully run %q", cmd.String())
	return nil
}

func CleanupNetorwk(name string) {
	log.Printf("TODO: really clean up things for %s", name)
}
