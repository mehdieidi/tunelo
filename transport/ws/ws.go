package ws

import (
	"nhooyr.io/websocket"

	"tunelo/pkg/logger"
	"tunelo/transport"
)

const ProtocolName = "ws"

type WebSocket struct {
	ServerAddr     string
	Logger         logger.Logger
	MsgHandlerFunc transport.MsgHandlerFunc
	Conn           *websocket.Conn
}

func New(
	serverAddr string,
	logger logger.Logger,
	msgHandlerFunc transport.MsgHandlerFunc,
) *WebSocket {
	return &WebSocket{
		ServerAddr:     serverAddr,
		Logger:         logger,
		MsgHandlerFunc: msgHandlerFunc,
	}
}
