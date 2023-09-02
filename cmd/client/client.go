package client

import (
	"fmt"
	"net"
	"os"

	"nhooyr.io/websocket"

	"tunelo/pkg/graceful"
	"tunelo/transport/udp"
	"tunelo/transport/ws"
	"tunelo/wire"
)

func Run(cfg Config) {
	clientAddr := net.JoinHostPort("127.0.0.1", cfg.ClientPort)
	serverAddr := net.JoinHostPort(cfg.ServerIP, cfg.ServerPort)
	wsAddr := fmt.Sprintf("ws://%s/ws", serverAddr)
	vpnAddr := net.JoinHostPort("127.0.0.1", cfg.VPNPort)

	udp := udp.New(
		clientAddr,
		cfg.Logger,
		nil,
		cfg.BufSize,
	)

	if err := udp.Listen(); err != nil {
		cfg.Logger.Error(err, nil)
	}
	defer udp.Listener.Close()

	cfg.Logger.Info("UDP server listening on "+clientAddr, nil)

	wsConn, err := ws.Dial(wsAddr)
	if err != nil {
		cfg.Logger.Error(fmt.Errorf("dialing ws: %v", err), nil)
		os.Exit(1)
	}
	defer wsConn.Close(websocket.StatusNormalClosure, "")

	cfg.Logger.Info("WebSocket connected: "+wsAddr, nil)

	websocket := ws.New(
		serverAddr,
		cfg.Logger,
		nil,
	)
	websocket.Conn = wsConn

	wire := wire.New(
		websocket,
		cfg.SecretKey,
		cfg.Logger,
		vpnAddr,
		cfg.BufSize,
	)

	websocket.MsgHandlerFunc = wire.WebSocketMsgHandler
	udp.MsgHandlerFunc = wire.UDPMsgHandler

	go websocket.Read()
	go udp.Read()

	graceful.ShutdownHandler()
}
