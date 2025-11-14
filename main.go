package main

import (
	"example.com/framespector/network"
	"log"
)

func main() {
	log.SetPrefix("framespector: ")
	log.SetFlags(0)

	if err := network.SetupNetwork("veth0"); err != nil {
		log.Fatal(err)
	}
	defer network.CleanupNetorwk("veth0")

	log.Println("Setup network done")
}
