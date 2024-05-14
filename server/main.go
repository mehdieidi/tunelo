package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"tunelo/pkg/logger/plain"
)

func main() {
	var serverIP string
	var serverPort string
	var vpnPort string
	var protocol string

	flag.StringVar(&serverIP, "server_ip", "127.0.0.1", "Proxy server IP address.")
	flag.StringVar(&serverPort, "server_port", "23230", "Proxy server port number.")
	flag.StringVar(&vpnPort, "vpn_port", "23233", "Local VPN port number.")
	flag.StringVar(&protocol, "p", "ws", "Tunnel transport protocol. Options: ws, utls, and tcp.")
	flag.Parse()

	logger := plain.New()

	vpnAddr := net.JoinHostPort("127.0.0.1", vpnPort)
	vpnUDPAddr, err := net.ResolveUDPAddr("udp", vpnAddr)
	if err != nil {
		logger.Error(fmt.Errorf("error resolving vpn udp addr: %v", err), nil)
		os.Exit(1)
	}

	localUDPAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:")
	if err != nil {
		logger.Error(fmt.Errorf("error resolving local udp addr: %v", err), nil)
		os.Exit(1)
	}

	vpnConn, err := net.DialUDP("udp", localUDPAddr, vpnUDPAddr)
	if err != nil {
		logger.Error(fmt.Errorf("error dialling vpn: %v", err), nil)
		os.Exit(1)
	}
	defer func(c *net.UDPConn) {
		err := c.Close()
		if err != nil {
			logger.Error(fmt.Errorf("error closing udp conn: %v", err), nil)
		}
	}(vpnConn)

	serverAddr := net.JoinHostPort(serverIP, serverPort)

	switch protocol {
	case "utls":
		t := utlsTransport{
			serverAddr: serverAddr,
			vpnConn:    vpnConn,
			logger:     logger,
		}
		err := t.run()
		if err != nil {
			logger.Error(err, nil)
			os.Exit(1)
		}
	case "tcp":
		t := tcpTransport{
			serverAddr: serverAddr,
			vpnConn:    vpnConn,
			logger:     logger,
		}
		err := t.run()
		if err != nil {
			logger.Error(err, nil)
			os.Exit(1)
		}
	default:
		t := wsTransport{
			serverAddr: serverAddr,
			vpnConn:    vpnConn,
			logger:     logger,
		}
		err := t.run()
		if err != nil {
			logger.Error(err, nil)
			os.Exit(1)
		}
	}
}
