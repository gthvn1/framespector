package network

import (
	"fmt"
)

func ProcessFrame(veth *Veth, data []byte) ([]byte, error) {
	f, err := parseEthernet(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	// Dispatch based on the ethernet type
	switch f.EtherType {
	case EtherTypeARP:
		handleARP(veth.PeerName, veth.PeerIP, f.Payload)
	case EtherTypeIPv4:
		handleIPv4(veth.PeerIP, f.Payload)
	case EtherTypeIPv6:
		return handleIPv6(f.Payload)
	case EtherTypeVLAN:
		return nil, fmt.Errorf("VLAN frame ignored")
	case EtherTypeUnknown:
		return nil, fmt.Errorf("unkown ether type %s", fmt.Sprintf("0x%04x", f.EtherType))
	default:
		// If you are here it is because you modified the EtherType enum and you
		// don't handle it here.
		panic(fmt.Sprintf("unhandled EtherType in switch: 0x%04x", f.EtherType))
	}
	return nil, nil
}
