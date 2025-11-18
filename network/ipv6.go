package network

func handleIPv6(payload []byte) ([]byte, error) {
	_ = payload
	return nil, &ToDoWarning{Msg: "handle IPv6 frame", EtherType: EtherTypeIPv6}
}
