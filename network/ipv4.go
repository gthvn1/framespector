package network

import (
	"encoding/binary"
	"fmt"
	"net"
)

// +--------------------------------------------------------+
// | IPv4 Header (20-60 bytes, typically 20)                |
// |--------------------------------------------------------|
// | Ver/IHL (1) | DSCP/ECN (1) | Total Length (2)          |
// | Identification (2) | Flags/Fragment Offset (2)         |
// | TTL (1) | Protocol (1) | Header Checksum (2)           |
// | Source IP (4) | Destination IP (4)                     |
// | Options (0-40 bytes, optional if IHL > 5)              |
// +--------------------------------------------------------+
//
// [RFC 791] https://datatracker.ietf.org/doc/html/rfc791
// Here is a raw ethernet frame that we received
//   -> ping from 192.168.38.2 to 192.168.38.3:
//
// ETH: 06 2b 41 e7 ae 3c
// ETH: 22 74 85 fe 7e 04
// ETH: 08 00  --> This is IP
// IP: 45            -> Version:4 (it is ipv4), Internet Header Length (IHL): 5
// IP: 00
// IP: 00 54         -> total length: 84 bytes (entire packet size in bytes, including header and data)
// IP: dd 7c         -> Identification: 22364
// IP: 40 00         -> Flags: Don't Fragment, Fragment Offset: 0
// IP: 40            -> TTL: 64
// IP: 01            -> Protocol: ICMP
// IP: 8f d6         -> Header Checksum: 36798
// IP: c0 a8 26 02   -> Source IP: 192.168.38.2
// IP: c0 a8 26 03   -> Destination IP: 192.168.38.3
// IP: 08 00 54 93 04 2b 00 01 29 f9
// IP: 11 69 00 00 00 00 a0 0b 05 00
// IP: 00 00 00 00 10 11 12 13 14 15
// IP: 16 17 18 19 1a 1b 1c 1d 1e 1f
// IP: 20 21 22 23 24 25 26 27 28 29
// IP: 2a 2b 2c 2d 2e 2f 30 31 32 33
// IP: 34 35 36 37

// https://en.wikipedia.org/wiki/IPv4#header
type IPv4Protocol = uint8

const (
	ICMPProtocol IPv4Protocol = 1
	TCPProtocol  IPv4Protocol = 6
	UDPProtocol  IPv4Protocol = 17
)

type IPv4Packet struct {
	// Version + IHL
	VersionIHL uint8
	// Differentiated Services Code Point + Explicit Congestion Notification
	DSCPECN uint8
	// Total packet lenght (header + data)
	TotalLength uint16
	// Fragment identification
	Identification uint16
	// Flags (bit 0: reserved, bit 1: DF, bit 2: MF) + Fragment offset
	FlagsFragOffset uint16
	// Time to live
	TTL uint8
	// Protocol (1=ICMP, 6=TCP, 17=UDP)
	Protocol IPv4Protocol
	// Header checksum
	HeaderChecksum uint16
	// Source IP address
	SourceIP net.IP
	// Destination IP address
	DestIP net.IP
	// Options (if IHL > 5)
	Options []byte
	// Payload data
	Payload []byte
}

func handleIPv4(peerIP net.IP, payload []byte) ([]byte, error) {
	p, err := parseIPv4Packet(payload, peerIP)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IPv4 packet: %w", err)
	}

	switch p.Protocol {
	case ICMPProtocol:
		icmp, err := parseICMP(p)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ICMP packet: %w", err)
		}

		if icmp.Type != ICMPEchoRequest {
			return nil, fmt.Errorf("only ICMP Echo request are handled")
		}

		return nil, fmt.Errorf("todo: handle ICMP echo request")
	default:
		return nil, fmt.Errorf("only ICMP protocol is managed currently")
	}
}

func parseIPv4Packet(payload []byte, ourIP net.IP) (*IPv4Packet, error) {
	if len(payload) < 20 {
		return nil, fmt.Errorf("IPv4 packet too short: %d bytes (minimum 20)", len(payload))
	}

	p := &IPv4Packet{
		VersionIHL:      payload[0],
		DSCPECN:         payload[1],
		TotalLength:     binary.BigEndian.Uint16(payload[2:4]),
		Identification:  binary.BigEndian.Uint16(payload[4:6]),
		FlagsFragOffset: binary.BigEndian.Uint16(payload[6:8]),
		TTL:             payload[8],
		Protocol:        IPv4Protocol(payload[9]),
		HeaderChecksum:  binary.BigEndian.Uint16(payload[10:12]),
		SourceIP:        net.IP(payload[12:16]),
		DestIP:          net.IP(payload[16:20]),
	}

	// Validate version
	if p.Version() != 4 {
		return nil, fmt.Errorf("not IPv4: version=%d", p.Version())
	}

	// Calculate header length
	headerLen := int(p.IHL()) * 4
	if headerLen < 20 {
		return nil, fmt.Errorf("invalid IHL: %d (too small)", p.IHL())
	}
	if headerLen > len(payload) {
		return nil, fmt.Errorf("IHL indicates %d bytes but packet is only %d bytes", headerLen, len(payload))
	}

	// Extract options if present (IHL > 5 means options exist)
	if headerLen > 20 {
		p.Options = make([]byte, headerLen-20)
		copy(p.Options, payload[20:headerLen])
	}

	// Extract payload
	if len(payload) > headerLen {
		p.Payload = payload[headerLen:]
	}

	return p, nil
}

// ------------------------------------------------------------------------------
// Accessor methods for packed fields
func (p *IPv4Packet) Version() uint8 {
	return (p.VersionIHL >> 4) & 0x0F
}

func (p *IPv4Packet) IHL() uint8 {
	return p.VersionIHL & 0x0F
}
