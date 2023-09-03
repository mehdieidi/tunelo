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

	wsTransport := &ws.WebSocket{
		ServerAddr: serverAddr,
		Logger:     cfg.Logger,
		Endpoint:   cfg.WebSocketEndpoint,
	}

	wire := &wire.Wire{
		SecretKey: cfg.SecretKey,
		Logger:    cfg.Logger,
		VPNAddr:   vpnAddr,
		BufSize:   cfg.BufSize,
	}

	wsTransport.MsgHandlerFunc = wire.WebSocketMsgHandler

	if err := wsTransport.ListenAndServe(); err != nil {
		cfg.Logger.Error(fmt.Errorf("websocket listen and serve: %v", err), nil)
		os.Exit(1)
	}
}
