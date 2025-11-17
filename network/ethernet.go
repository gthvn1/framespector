package network

import (
	"encoding/binary"
	"fmt"
	"net"
)

// +--------------------------------------------------------+
// | Ethernet Header (14 bytes standard)                    |
// |--------------------------------------------------------|
// | Destination MAC (6) | Source MAC (6) | EtherType (2)   |
// +--------------------------------------------------------+
//
// Ethernet II layout begins with:
//
//	Destination MAC: 6 bytes
//	Source MAC: 6 bytes
//	EtherType: 2 bytes (0x0806 -> ARP)
//
// https://en.wikipedia.org/wiki/Ethernet_frame
// https://en.wikipedia.org/wiki/EtherType
type EtherType uint16

const (
	EtherTypeIPv4    EtherType = 0x0800
	EtherTypeARP     EtherType = 0x0806
	EtherTypeIPv6    EtherType = 0x86DD
	EtherTypeVLAN    EtherType = 0x8100
	EtherTypeUnknown EtherType = 0xFFFF
)

func parseEtherType(v uint16) EtherType {
	// Default to unknown
	switch EtherType(v) {
	case EtherTypeARP, EtherTypeIPv4, EtherTypeIPv6, EtherTypeVLAN:
		return EtherType(v)
	default:
		return EtherTypeUnknown
	}
}

func (e EtherType) string() string {
	switch e {
	case EtherTypeIPv4:
		return "IPv4"
	case EtherTypeARP:
		return "ARP"
	case EtherTypeIPv6:
		return "IPv6"
	case EtherTypeVLAN:
		return "VLAN"
	case EtherTypeUnknown:
		return "Unknown"
	default:
		return fmt.Sprintf("0x%04X", uint16(e))
	}
}

type EthernetFrame struct {
	DestMAC   net.HardwareAddr
	SrcMAC    net.HardwareAddr
	EtherType EtherType
	HeaderLen int
	Payload   []byte
}

func handleARP(peerName string, peerIP net.IP, payload []byte) ([]byte, error) {
	peerIface, err1 := net.InterfaceByName(peerName)
	if err1 != nil {
		return nil, fmt.Errorf("failed to peer interface %s: %w", peerName, err1)
	}

	reply, err2 := parseARP(payload, peerIface.HardwareAddr, peerIP)
	if err2 != nil {
		return nil, fmt.Errorf("ARP request not handled: %w", err2)
	}

	arpPayload := reply.marshal()
	return buildEthernetFrame(reply.TargetHA, reply.SenderHA, EtherTypeARP, arpPayload), nil
}

func parseEthernet(packet []byte) (*EthernetFrame, error) {
	if len(packet) < 14 {
		return nil, fmt.Errorf("packet too small: need at least 14 bytes, got %d", len(packet))
	}

	f := &EthernetFrame{
		DestMAC: net.HardwareAddr(packet[0:6]),
		SrcMAC:  net.HardwareAddr(packet[6:12]),
	}

	// At offset we need to read etherType and check if it is vlan and
	// handle VLAN tag (802.1Q)
	offset := 12
	et := binary.BigEndian.Uint16(packet[offset : offset+2])
	if et == uint16(EtherTypeVLAN) {
		offset += 4 // Skip 4-byte VLAN tag
		if len(packet) < offset+2 {
			return nil, fmt.Errorf("packet too small for VLAN: need at least %d bytes", offset+2)
		}
		et = binary.BigEndian.Uint16(packet[offset : offset+2])
	}

	f.EtherType = parseEtherType(et)
	f.HeaderLen = offset + 2
	f.Payload = packet[f.HeaderLen:]

	// For debugging purpose print raw ARP frame
	if f.EtherType == EtherTypeARP {
		fmt.Println("--------- ARP FRAME ---------")
		printHex(packet)
		fmt.Println("-----------------------------")
	}

	return f, nil
}

func buildEthernetFrame(dst, src net.HardwareAddr, etherType EtherType, payload []byte) []byte {
	frame := make([]byte, 14+len(payload))

	// MAC addresses
	copy(frame[0:6], dst)  // destination
	copy(frame[6:12], src) // source

	// EtherType in big-endian
	binary.BigEndian.PutUint16(frame[12:14], uint16(etherType))

	// Payload
	copy(frame[14:], payload)

	return frame
}

// string returns a human-readable representation
//func (f *EthernetFrame) string() string {
//	return fmt.Sprintf("Ethernet: %s -> %s, Type: %s, Payload: %d bytes",
//		f.SrcMAC, f.DestMAC, f.EtherType.string(), len(f.Payload))
//}

func printHex(buf []byte) {
	for i := 0; i < len(buf); i += 10 {
		end := min(i+10, len(buf))

		for _, b := range buf[i:end] {
			fmt.Printf("%02x ", b)
		}

		fmt.Println()
	}
}
