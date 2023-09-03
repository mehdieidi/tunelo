package client

import (
	"fmt"
	"net"
	"os"

	"nhooyr.io/websocket"

	"tunelo/pkg/graceful"
	"tunelo/transport"
	"tunelo/transport/udp"
	"tunelo/transport/ws"
	"tunelo/wire"
)

func Run(cfg Config) {
	clientAddr := net.JoinHostPort("127.0.0.1", cfg.ClientPort)
	serverAddr := net.JoinHostPort(cfg.ServerIP, cfg.ServerPort)
	wsAddr := fmt.Sprintf("ws://%s/%s", serverAddr, cfg.WebSocketEndpoint)
	vpnAddr := net.JoinHostPort("127.0.0.1", cfg.VPNPort)

	udp := &udp.UDP{
		ServerAddr: clientAddr,
		Logger:     cfg.Logger,
		BufSize:    cfg.BufSize,
	}

	if err := udp.Listen(); err != nil {
		cfg.Logger.Error(fmt.Errorf("creating udp listener: %v", err), nil)
		os.Exit(1)
	}
	defer udp.Listener.Close()

	cfg.Logger.Info(fmt.Sprintf("UDP server listening on %s", clientAddr), nil)

	wsConn, err := ws.Dial(wsAddr)
	if err != nil {
		cfg.Logger.Error(fmt.Errorf("dialing ws: %v", err), nil)
		os.Exit(1)
	}
	defer wsConn.Close(websocket.StatusNormalClosure, "")

	transportConn := &transport.Conn{WebSocket: wsConn}

	cfg.Logger.Info(fmt.Sprintf("WebSocket connected: %s", wsAddr), nil)

	wsTransport := &ws.WebSocket{
		ServerAddr: serverAddr,
		Logger:     cfg.Logger,
	}

	wire := &wire.Wire{
		SecretKey: cfg.SecretKey,
		Logger:    cfg.Logger,
		VPNAddr:   vpnAddr,
		BufSize:   cfg.BufSize,
	}

	wsTransport.MsgHandlerFunc = wire.WebSocketMsgHandler
	udp.MsgHandlerFunc = wire.UDPMsgHandler

	go wsTransport.Read(transportConn)
	go udp.Read(transportConn)

	graceful.ShutdownHandler()
}
