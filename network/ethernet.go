package network

import "fmt"

type EthernetFrame struct {
	DestMacAddr [6]byte
	SrcMacAddr  [6]byte
	EtherType   uint16
	Payload     []byte
}

func ParseEthernet(data []byte) (*EthernetFrame, error) {
	fmt.Println("--------- RAW FRAME ---------")
	printHex(data)
	fmt.Println("-----------------------------")
	return nil, nil
}

func printHex(buf []byte) {
	for i := 0; i < len(buf); i += 10 {
		end := min(i+10, len(buf))

		for _, b := range buf[i:end] {
			fmt.Printf("%02x ", b)
		}

		fmt.Println()
	}
}
