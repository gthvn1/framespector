package network

import (
	"errors"
	"fmt"
)

var DecodeDataError = errors.New("failed to decode data")

type ToDoWarning struct {
	Msg       string
	EtherType EtherType
}

func (e *ToDoWarning) Error() string {
	return fmt.Sprintf("todo: %s for %s", e.Msg, e.EtherType.String())
}

func ProcessFrame(veth *Veth, data []byte) ([]byte, error) {
	f, err := parseEthernet(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", DecodeDataError, err)
	}

	// Dispatch based on the ethernet type
	switch f.EtherType {
	case EtherTypeARP:
		return handleARP(veth.PeerName, veth.PeerIP, f.Payload)
	case EtherTypeIPv4:
		return handleIPv4(veth.PeerIP, f.Payload)
	case EtherTypeIPv6:
		return handleIPv6(f.Payload)
	case EtherTypeVLAN, EtherTypeUnknown:
		return nil, &ToDoWarning{Msg: "should we handle this", EtherType: f.EtherType}
	default:
		// If you are here it is because you modified the EtherType enum and you
		// don't handle it here.
		panic(fmt.Sprintf("unhandled EtherType in switch: %s", f.EtherType.String()))
	}
}
