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
	BufSize        int
}

func New(
	serverAddr string,
	logger logger.Logger,
	msgHandlerFunc MsgHandlerFunc,
	bufSize int,
) *UDP {
	return &UDP{
		ServerAddr:     serverAddr,
		Logger:         logger,
		MsgHandlerFunc: msgHandlerFunc,
		BufSize:        bufSize,
	}
}
