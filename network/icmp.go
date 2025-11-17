package network

import (
	"encoding/binary"
	"fmt"
)

// https://datatracker.ietf.org/doc/html/rfc792
// https://en.wikipedia.org/wiki/Internet_Control_Message_Protocol
//
// Here is an example of raw payload:
// 08          -> Type
// 00          -> Code
// 12 76       -> checksum
// 0a e3 00 01 -> Rest of the header
// db 10 12 69 00 00 00 00 29 59
// 05 00 00 00 00 00 10 11 12 13
// 14 15 16 17 18 19 1a 1b 1c 1d
// 1e 1f 20 21 22 23 24 25 26 27
// 28 29 2a 2b 2c 2d 2e 2f 30 31
// 32 33 34 35 36 37
type ICMPType = uint8

// Focusing on responding to ping
const (
	ICMPEchoReply   ICMPType = 0
	ICMPEchoRequest ICMPType = 8
)

type ICMPPacket struct {
	Type           ICMPType
	Code           uint8
	Checksum       uint16
	Identifier     uint16
	SequenceNumber uint16
	Data           []byte
}

func ParseICMP(packet *IPv4Packet) (*ICMPPacket, error) {
	// Minimal size is 8 bytes
	payload := packet.Payload

	if len(payload) < 8 {
		return nil, fmt.Errorf("ICMP packet too short")
	}

	p := &ICMPPacket{
		Type:           ICMPType(payload[0]),
		Code:           payload[1],
		Checksum:       binary.BigEndian.Uint16(payload[2:4]),
		Identifier:     binary.BigEndian.Uint16(payload[4:6]),
		SequenceNumber: binary.BigEndian.Uint16(payload[6:8]),
		Data:           payload[8:],
	}

	return p, nil
}

func (p *ICMPPacket) Marshal() []byte {
	data := make([]byte, 8+len(p.Data))

	data[0] = byte(p.Type)
	data[1] = p.Code
	binary.BigEndian.PutUint16(data[4:6], p.Identifier)
	binary.BigEndian.PutUint16(data[6:8], p.SequenceNumber)
	copy(data[8:], p.Data)

	// compute checksum
	binary.BigEndian.PutUint16(data[2:4], 0)
	cs := checksum(data)
	binary.BigEndian.PutUint16(data[2:4], cs)

	return data
}

func checksum(data []byte) uint16 {
	var sum uint32
	n := len(data)

	for i := 0; i < n-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i:]))
	}

	if n%2 == 1 {
		sum += uint32(data[n-1]) << 8
	}

	for (sum >> 16) > 0 {
		sum = (sum >> 16) + (sum & 0xFFFF)
	}

	return ^uint16(sum)
}
