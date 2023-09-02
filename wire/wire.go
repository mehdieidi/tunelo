package wire

import (
	"tunelo/pkg/logger"
	"tunelo/transport/ws"
)

type Wire struct {
	WebSocket *ws.WebSocket
	SecretKey []byte
	Logger    logger.Logger
	VPNAddr   string
	BufSize   int
}

func New(
	websocket *ws.WebSocket,
	secretKey []byte,
	logger logger.Logger,
	vpnAddr string,
	bufSize int,
) *Wire {
	return &Wire{
		WebSocket: websocket,
		SecretKey: secretKey,
		Logger:    logger,
		VPNAddr:   vpnAddr,
		BufSize:   bufSize,
	}
}
