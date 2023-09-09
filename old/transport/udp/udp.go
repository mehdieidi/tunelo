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
