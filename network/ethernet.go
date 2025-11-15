package network

import "fmt"

type EthernetFrame struct {
	DestMacAddr [6]byte
	SrcMacAddr  [6]byte
	EtherType   uint16
	Payload     []byte
}

func ParseEthernet(data []byte) (*EthernetFrame, error) {
	fmt.Printf("Received %s", data)
	return nil, nil
}
