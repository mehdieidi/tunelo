package server

import (
	"fmt"
	"net"
	"os"

	"tunelo/transport/ws"
	"tunelo/wire"
)

func Run(cfg Config) {
	vpnAddr := net.JoinHostPort("127.0.0.1", cfg.VPNPort)
	serverAddr := net.JoinHostPort(cfg.ServerIP, cfg.ServerPort)

	ws := ws.New(serverAddr, cfg.Logger, nil)

	wire := wire.New(
		ws,
		cfg.SecretKey,
		cfg.Logger,
		vpnAddr,
		cfg.BufSize,
	)

	ws.MsgHandlerFunc = wire.WebSocketMsgHandler

	if err := wire.WebSocket.ListenAndServe(); err != nil {
		cfg.Logger.Error(fmt.Errorf("websocket listen and serve: %v", err), nil)
		os.Exit(1)
	}
}
