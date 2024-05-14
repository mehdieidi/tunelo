package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"tunelo/pkg/logger/plain"
)

func main() {
	var serverIP string
	var serverPort string
	var vpnPort string
	var clientPort string
	var protocol string
	var serverDomain string

	flag.StringVar(&serverIP, "server_ip", "127.0.0.1", "Remote proxy-server IP address.")
	flag.StringVar(&serverPort, "server_port", "23230", "Remote proxy-server port number.")
	flag.StringVar(&vpnPort, "vpn_port", "23233", "Local VPN port number.")
	flag.StringVar(&clientPort, "client_port", "23231", "Client port number.")
	flag.StringVar(&protocol, "p", "ws", "Tunnel transport protocol. Options: ws, utls, and tcp.")
	flag.StringVar(&serverDomain, "server_domain", "", "Server domain.")
	flag.Parse()

	logger := plain.New()

	clientAddr := net.JoinHostPort("127.0.0.1", clientPort)
	clientUDPAddr, err := net.ResolveUDPAddr("udp", clientAddr)
	if err != nil {
		logger.Error(fmt.Errorf("error resolving client udp addr: %v", err), nil)
		os.Exit(1)
	}

	clientUDPConn, err := net.ListenUDP("udp", clientUDPAddr)
	if err != nil {
		logger.Error(fmt.Errorf("error creating client udp listener: %v", err), nil)
		os.Exit(1)
	}
	defer func(c *net.UDPConn) {
		err := c.Close()
		if err != nil {
			logger.Error(fmt.Errorf("error closing client udp connection: %v", err), nil)
		}
	}(clientUDPConn)

	logger.Info(fmt.Sprintf("UDP conn: %s", clientAddr), nil)

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
			logger.Error(fmt.Errorf("error closing vpn connection: %v", err), nil)
		}
	}(vpnConn)

	logger.Info(fmt.Sprintf("VPN UDP conn: %s", vpnAddr), nil)

	serverAddr := net.JoinHostPort(serverIP, serverPort)

	switch protocol {
	case "utls":
		t := utlsTransport{
			serverDomain:  serverDomain,
			serverAddr:    serverAddr,
			vpnConn:       vpnConn,
			clientUDPConn: clientUDPConn,
			logger:        logger,
		}
		err := t.run()
		if err != nil {
			logger.Error(err, nil)
			os.Exit(1)
		}
	case "tcp":
		t := tcpTransport{
			serverAddr:    serverAddr,
			vpnConn:       vpnConn,
			clientUDPConn: clientUDPConn,
			logger:        logger,
		}
		err := t.run()
		if err != nil {
			logger.Error(err, nil)
			os.Exit(1)
		}
	default:
		t := wsTransport{
			serverAddr:    serverAddr,
			vpnConn:       vpnConn,
			clientUDPConn: clientUDPConn,
			logger:        logger,
		}
		err := t.run()
		if err != nil {
			logger.Error(err, nil)
			os.Exit(1)
		}
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	<-ctx.Done()

	fmt.Println("\n[-] shutdown signal received")
	os.Exit(0)
}
