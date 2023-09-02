package ws

import (
	"nhooyr.io/websocket"

	"tunelo/pkg/logger"
)

const ProtocolName = "ws"

type WebSocket struct {
	ServerAddr     string
	Logger         logger.Logger
	MsgHandlerFunc MsgHandlerFunc
	Conn           *websocket.Conn
}

func New(serverAddr string, logger logger.Logger, msgHandlerFunc MsgHandlerFunc) *WebSocket {
	return &WebSocket{
		ServerAddr:     serverAddr,
		Logger:         logger,
		MsgHandlerFunc: msgHandlerFunc,
	}
}
