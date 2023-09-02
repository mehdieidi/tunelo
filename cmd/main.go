package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"tunelo/cmd/client"
	"tunelo/cmd/server"
	"tunelo/pkg/logger/zerolog"
)

func main() {
	var vpnPort string
	var clientPort string
	var remoteServerIP string
	var remoteServerPort string
	var serverIP string
	var serverPort string
	var clientMode bool

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
	flag.BoolVar(
		&clientMode,
		"c",
		false,
		"Run in client mode. The default is server mode.",
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

	if clientMode {
		cfg := client.Config{
			VPNPort:          vpnPort,
			ClientPort:       clientPort,
			RemoteServerIP:   remoteServerIP,
			RemoteServerPort: remoteServerPort,
			SecretKey:        secretKey,
			Logger:           logger,
		}
		client.Run(cfg)
	} else {
		cfg := server.Config{
			VPNPort:    vpnPort,
			ServerIP:   serverIP,
			ServerPort: serverPort,
			SecretKey:  secretKey,
			Logger:     logger,
		}
		server.Run(cfg)
	}
}
