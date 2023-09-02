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
	var remoteServerIP string
	var remoteServerPort string
	var serverIP string
	var serverPort string
	var bufSize int
	var serverMode bool

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
		"Remote server IP used by the client to tunnel data.",
	)
	flag.StringVar(
		&remoteServerPort,
		"server_port",
		"23230",
		"Remote server port used by the client to tunnel data.",
	)
	flag.StringVar(
		&serverIP,
		"i",
		"127.0.0.1",
		"Server IP address of the tunnel.",
	)
	flag.StringVar(
		&serverPort,
		"p",
		"23230",
		"Server port number that the tunnel server will listen to.",
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

	flag.Parse()

	logger := plain.New(os.Stdout)

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
			VPNPort:    vpnPort,
			ServerIP:   serverIP,
			ServerPort: serverPort,
			SecretKey:  secretKey,
			Logger:     logger,
			BufSize:    bufSize,
		}
		server.Run(cfg)
	} else {
		cfg := client.Config{
			VPNPort:          vpnPort,
			ClientPort:       clientPort,
			RemoteServerIP:   remoteServerIP,
			RemoteServerPort: remoteServerPort,
			SecretKey:        secretKey,
			Logger:           logger,
			BufSize:          bufSize,
		}
		client.Run(cfg)
	}
}
