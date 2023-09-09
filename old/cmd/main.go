package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"tunelo/cmd/client"
	"tunelo/cmd/server"
	"tunelo/pkg/logger/plain"
)

func main() {
	var vpnPort string
	var clientPort string
	var serverIP string
	var serverPort string
	var bufSize int
	var serverMode bool
	var websocketEndpoint string

	flag.StringVar(
		&vpnPort,
		"vpn_port",
		"23233",
		"Port number that the VPN (e.g. WireGuard) is listening to.",
	)
	flag.StringVar(
		&clientPort,
		"client_port",
		"23231",
		"Port number that the client listens to.",
	)
	flag.StringVar(
		&serverIP,
		"server_ip",
		"127.0.0.1",
		"Server IP of the tunnel.",
	)
	flag.StringVar(
		&serverPort,
		"server_port",
		"23230",
		"Server port of the tunnel.",
	)
	flag.IntVar(
		&bufSize,
		"buf",
		1450,
		"Buffer size.",
	)
	flag.BoolVar(
		&serverMode,
		"s",
		false,
		"Run in server mode. The default is client mode.",
	)
	flag.StringVar(
		&websocketEndpoint,
		"ws_endpoint",
		"ws",
		"WebSocket endpoint that accepts websocket connections.",
	)
	flag.Parse()

	logger := plain.New()

	if err := godotenv.Load(); err != nil {
		logger.Error(fmt.Errorf("loading env: %v", err), nil)
		os.Exit(1)
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if string(secretKey) == "" {
		logger.Error(fmt.Errorf("secret key cannot be empty"), nil)
		os.Exit(1)
	}

	if serverMode {
		cfg := server.Config{
			VPNPort:           vpnPort,
			ServerIP:          serverIP,
			ServerPort:        serverPort,
			SecretKey:         secretKey,
			Logger:            logger,
			BufSize:           bufSize,
			WebSocketEndpoint: websocketEndpoint,
		}
		server.Run(cfg)
	} else {
		cfg := client.Config{
			VPNPort:           vpnPort,
			ClientPort:        clientPort,
			ServerIP:          serverIP,
			ServerPort:        serverPort,
			SecretKey:         secretKey,
			Logger:            logger,
			BufSize:           bufSize,
			WebSocketEndpoint: websocketEndpoint,
		}
		client.Run(cfg)
	}
}
