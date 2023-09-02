package udp

import (
	"net"

	"tunelo/pkg/logger"
)

type UDP struct {
	ServerAddr     string
	Logger         logger.Logger
	MsgHandlerFunc MsgHandlerFunc
	Conn           net.PacketConn
}

func New(
	serverAddr string,
	logger logger.Logger,
	msgHandlerFunc MsgHandlerFunc,
) *UDP {
	return &UDP{
		ServerAddr:     serverAddr,
		Logger:         logger,
		MsgHandlerFunc: msgHandlerFunc,
	}
}
