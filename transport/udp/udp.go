package udp

import (
	"net"

	"tunelo/pkg/logger"
	"tunelo/transport"
)

const ProtocolName = "udp"

type UDP struct {
	ServerAddr     string
	Logger         logger.Logger
	MsgHandlerFunc transport.MsgHandlerFunc
	Listener       net.PacketConn
	BufSize        int
}

func New(
	serverAddr string,
	logger logger.Logger,
	msgHandlerFunc transport.MsgHandlerFunc,
	bufSize int,
) *UDP {
	return &UDP{
		ServerAddr:     serverAddr,
		Logger:         logger,
		MsgHandlerFunc: msgHandlerFunc,
		BufSize:        bufSize,
	}
}
