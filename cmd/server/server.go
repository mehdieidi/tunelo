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

	websocket := ws.New(serverAddr, cfg.Logger, nil)

	wire := wire.New(
		websocket,
		cfg.SecretKey,
		cfg.Logger,
		vpnAddr,
		cfg.BufSize,
	)

	websocket.MsgHandlerFunc = wire.WebSocketMsgHandler

	if err := wire.WebSocket.ListenAndServe(); err != nil {
		cfg.Logger.Error(fmt.Errorf("websocket listen and serve: %v", err), nil)
		os.Exit(1)
	}
}
