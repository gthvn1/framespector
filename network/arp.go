package network

import (
	"log/slog"
)

func HandleARP(logger *slog.Logger, payload []byte) {
	_ = payload
	logger.Warn("TODO: handle ARP")
}
