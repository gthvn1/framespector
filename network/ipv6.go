package network

import (
	"fmt"
)

func handleIPv6(payload []byte) ([]byte, error) {
	_ = payload
	return nil, fmt.Errorf("todo: decode ipv6")
}
