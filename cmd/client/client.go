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
	serverAddr := net.JoinHostPort(cfg.RemoteServerIP, cfg.RemoteServerPort)
	wsServerAddr := fmt.Sprintf("ws://%s/ws", serverAddr)
	vpnAddr := net.JoinHostPort("127.0.0.1", cfg.VPNPort)

	udp := udp.New(
		clientAddr,
		cfg.Logger,
		nil,
	)

	if err := udp.Listen(); err != nil {
		cfg.Logger.Error(fmt.Errorf("[error] creating udp listener: %v", err), nil)
	}
	defer udp.Conn.Close()

	cfg.Logger.Info("[info] UDP server listening to "+clientAddr, nil)

	wsConn, err := ws.Dial(wsServerAddr)
	if err != nil {
		cfg.Logger.Error(fmt.Errorf("[error] dialing ws: %v", err), nil)
		os.Exit(1)
	}
	defer wsConn.Close(websocket.StatusInternalError, "")

	cfg.Logger.Info("[info] WebSocket connected: "+wsServerAddr, nil)

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
		1450,
	)

	websocket.MsgHandlerFunc = wire.WebSocketMsgHandler
	udp.MsgHandlerFunc = wire.UDPMsgHandler

	go websocket.Read()
	go udp.UDPReadHandler()

	graceful.ShutdownHandler()
}
