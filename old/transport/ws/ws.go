package ws

import (
	"tunelo/pkg/logger"
	"tunelo/transport"
)

const ProtocolName = "ws"

type WebSocket struct {
	ServerAddr     string
	Logger         logger.Logger
	MsgHandlerFunc transport.MsgHandlerFunc
	Endpoint       string
}
