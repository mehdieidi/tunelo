package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
	"nhooyr.io/websocket"

	"tunelo/pkg/graceful"
	"tunelo/pkg/logger/zerolog"
	"tunelo/transport/udp"
	"tunelo/transport/ws"
	"tunelo/wire"
)

func main() {
	var vpnPort string
	var clientPort string
	var remoteServerIP string
	var remoteServerPort string

	flag.StringVar(
		&vpnPort,
		"vpn_port",
		"23233",
		"Port number that the VPN (e.g. WireGuard) listens to.",
	)
	flag.StringVar(
		&clientPort,
		"l",
		"23231",
		"Port number that the client listens to.",
	)
	flag.StringVar(
		&remoteServerIP,
		"server_ip",
		"127.0.0.1",
		"Remote server IP.",
	)
	flag.StringVar(
		&remoteServerPort,
		"server_port",
		"23230",
		"Remote server port.",
	)
	flag.Parse()

	logger := zerolog.New(os.Stdout)

	if err := godotenv.Load(); err != nil {
		logger.Error(fmt.Errorf("[error] loading env: %v", err), nil)
		os.Exit(1)
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if string(secretKey) == "" {
		logger.Error(fmt.Errorf("[error] secret key cannot be empty"), nil)
		os.Exit(1)
	}

	clientAddr := net.JoinHostPort("127.0.0.1", clientPort)
	serverAddr := net.JoinHostPort(remoteServerIP, remoteServerPort)
	wsServerAddr := fmt.Sprintf("ws://%s/ws", serverAddr)
	vpnAddr := net.JoinHostPort("127.0.0.1", vpnPort)

	udp := udp.New(
		clientAddr,
		logger,
		nil,
	)

	if err := udp.Listen(); err != nil {
		logger.Error(fmt.Errorf("[error] creating udp listener: %v", err), nil)
	}
	defer udp.Conn.Close()

	logger.Info("[info] UDP server listening to "+clientAddr, nil)

	wsConn, err := ws.Dial(wsServerAddr)
	if err != nil {
		logger.Error(fmt.Errorf("[error] dialing ws: %v", err), nil)
		os.Exit(1)
	}
	defer wsConn.Close(websocket.StatusInternalError, "")

	logger.Info("[info] WebSocket connected: "+wsServerAddr, nil)

	websocket := ws.New(
		serverAddr,
		logger,
		nil,
	)
	websocket.Conn = wsConn

	wire := wire.New(
		websocket,
		secretKey,
		logger,
		vpnAddr,
		1450,
	)

	websocket.MsgHandlerFunc = wire.WebSocketMsgHandler
	udp.MsgHandlerFunc = wire.UDPMsgHandler

	go websocket.Read()
	go udp.UDPReadHandler()

	graceful.ShutdownHandler()
}
