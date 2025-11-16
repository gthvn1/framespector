package network

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
)

// +--------------------------------------------------------+
// | ARP Payload (28 bytes standard for Ethernet/IPv4)      |
// |--------------------------------------------------------|
// | HTYPE (2) | PTYPE (2) | HLEN (1) | PLEN (1) | OPER (2) |
// | SHA (6) | SPA (4) | THA (6) | TPA (4)                  |
// +--------------------------------------------------------+
//
// [RFC ARP] https://datatracker.ietf.org/doc/html/rfc826
//
// Here is an example of what we are receiving from arping:
// ff ff ff ff ff ff  -> ETHERNET: broadcast
// f2 4e 68 82 e2 1b  -> ETHERNET: sender's MAC address
// 08 06              -> ETHERNET: ARP protocol
//
// 00 01              -> ARP: Hardware type
// 08 00              -> ARP: Protocl type
// 06                 -> ARP: Hardware size
// 04                 -> ARP: Protocol size
// 00 01              -> ARP: Opcode
// f2 4e 68 82 e2 1b  -> ARP: Sender MAC
// c0 a8 26 02        -> ARP: Sender IP
// ff ff ff ff ff ff  -> ARP: Target MAC
// c0 a8 26 03        -> ARP: Target IP
//
// https://en.wikipedia.org/wiki/Address_Resolution_Protocol
type ARPOper uint16

const (
	ARPRequest ARPOper = 1
	ARPReply   ARPOper = 2
)

type ARPPacket struct {
	HWType uint16 // Hardware type (1 = Ethernet)
	PType  uint16 // Protocol type
	HWLen  uint8  // Hardware address length (6 for MAC)
	PLen   uint8  // Protocol address length (4 for IPv4)
	Oper   ARPOper

	SenderHA net.HardwareAddr // Sender hardware address
	SenderPA net.IP           // Sender protocol address
	TargetHA net.HardwareAddr // Target hardware address
	TargetPA net.IP           // Target protocol address
}

func HandleARP(logger *slog.Logger, payload []byte, ourMAC net.HardwareAddr, ourIP net.IP) (*ARPPacket, error) {
	p, err := parseARP(payload)
	if err != nil {
		return nil, err
	}

	if p.Oper != ARPRequest {
		return nil, fmt.Errorf("only answer to ARP request")
	}

	if !p.TargetPA.Equal(ourIP) {
		return nil, fmt.Errorf("IP %s is not matching %s", ourIP.String(), p.TargetPA.String())
	}

	reply := &ARPPacket{
		HWType:   p.HWType,
		PType:    p.PType,
		HWLen:    p.HWLen,
		PLen:     p.PLen,
		Oper:     ARPReply,
		SenderHA: ourMAC,
		SenderPA: ourIP,
		TargetHA: p.SenderHA,
		TargetPA: p.SenderPA,
	}

	return reply, nil
}

func (p *ARPPacket) Marshal() []byte {
	b := make([]byte, 8+int(p.HWLen)*2+int(p.PLen)*2)

	binary.BigEndian.PutUint16(b[0:2], p.HWType)
	binary.BigEndian.PutUint16(b[2:4], p.PType)
	b[4] = p.HWLen
	b[5] = p.PLen
	binary.BigEndian.PutUint16(b[6:8], uint16(p.Oper))

	offset := 8
	copy(b[offset:offset+int(p.HWLen)], p.SenderHA)
	offset += int(p.HWLen)

	senderIP := p.SenderPA.To4()
	if senderIP == nil {
		senderIP = p.SenderPA // fallback if not IPv4
	}
	copy(b[offset:offset+int(p.PLen)], senderIP)
	offset += int(p.PLen)

	copy(b[offset:offset+int(p.HWLen)], p.TargetHA)
	offset += int(p.HWLen)

	targetIP := p.TargetPA.To4()
	if targetIP == nil {
		targetIP = p.TargetPA // fallback if not IPv4
	}

	copy(b[offset:offset+int(p.PLen)], targetIP)

	return b
}

func parseARP(payload []byte) (*ARPPacket, error) {
	// To get the operation we need at least 8 bytes
	if len(payload) < 8 {
		return nil, fmt.Errorf("ARP packet too small: need at least 8 bytes, got %d", len(payload))
	}

	p := &ARPPacket{
		HWType: binary.BigEndian.Uint16(payload[0:2]),
		PType:  binary.BigEndian.Uint16(payload[2:4]),
		HWLen:  payload[4],
		PLen:   payload[5],
		Oper:   ARPOper(binary.BigEndian.Uint16(payload[6:8])),
	}

	// Now we can compute the expected len
	expected := 8 + int(2*p.HWLen) + int(2*p.PLen)
	if len(payload) < expected {
		return nil, fmt.Errorf("ARP packet invalid len: expected %d bytes, got %d", expected, len(payload))
	}

	// Offsets for variable fields
	offset := 8

	end := offset + int(p.HWLen)
	p.SenderHA = net.HardwareAddr(payload[offset:end])
	offset = end

	end = offset + int(p.PLen)
	p.SenderPA = net.IP(payload[offset:end])
	offset = end

	end = offset + int(p.HWLen)
	p.TargetHA = net.HardwareAddr(payload[offset:end])
	offset = end

	end = offset + int(p.PLen)
	p.TargetPA = net.IP(payload[offset:end])

	return p, nil
}
