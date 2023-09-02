package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"

	"tunelo/pkg/logger/zerolog"
	"tunelo/transport/ws"
	"tunelo/wire"
)

func main() {
	var vpnPort string
	var serverIP string
	var serverPort string

	flag.StringVar(
		&vpnPort,
		"vpn_port",
		"23233",
		"Port number that the VPN (e.g. WireGuard) listens to.",
	)
	flag.StringVar(
		&serverIP,
		"i",
		"127.0.0.1",
		"Server IP address.",
	)
	flag.StringVar(
		&serverPort,
		"p",
		"23230",
		"Server Port number.",
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

	vpnAddr := net.JoinHostPort("127.0.0.1", vpnPort)
	serverAddr := net.JoinHostPort(serverIP, serverPort)

	ws := ws.New(serverAddr, logger, nil)

	wire := wire.New(
		ws,
		secretKey,
		logger,
		vpnAddr,
		1450,
	)

	ws.MsgHandlerFunc = wire.WebSocketMsgHandler

	if err := wire.WebSocket.ListenAndServe(); err != nil {
		logger.Error(fmt.Errorf("[error] websocket listen and serve: %v", err), nil)
		os.Exit(1)
	}
}
