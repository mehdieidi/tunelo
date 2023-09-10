package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"nhooyr.io/websocket"

	"tunelo/pkg/logger/plain"
)

func main() {
	var serverIP string
	var serverPort string
	var vpnPort string
	var protocol string

	flag.StringVar(&serverIP, "server_ip", "127.0.0.1", "Remote proxy-server IP address.")
	flag.StringVar(&serverPort, "server_port", "23230", "Remote proxy-server port number.")
	flag.StringVar(&vpnPort, "vpn_port", "23233", "Local VPN port number.")
	flag.StringVar(&protocol, "p", "ws", "Tunnel transport protocol. Options: ws and tcp.")
	flag.Parse()

	log := plain.New()

	clientAddr := net.JoinHostPort("127.0.0.1", "23231")
	clientUDPAddr, err := net.ResolveUDPAddr("udp", clientAddr)
	if err != nil {
		log.Error(fmt.Errorf("resolving client udp addr: %v", err), nil)
		os.Exit(1)
	}

	clientUDPConn, err := net.ListenUDP("udp", clientUDPAddr)
	if err != nil {
		log.Error(fmt.Errorf("creating client udp listener: %v", err), nil)
		os.Exit(1)
	}
	defer clientUDPConn.Close()

	log.Info(fmt.Sprintf("UDP conn: %s", clientAddr), nil)

	vpnAddr := net.JoinHostPort("127.0.0.1", vpnPort)
	vpnUDPAddr, err := net.ResolveUDPAddr("udp", vpnAddr)
	if err != nil {
		log.Error(fmt.Errorf("resolving vpn udp addr: %v", err), nil)
		os.Exit(1)
	}

	localUDPAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:")
	if err != nil {
		log.Error(fmt.Errorf("resolving local udp addr: %v", err), nil)
		os.Exit(1)
	}

	vpnConn, err := net.DialUDP("udp", localUDPAddr, vpnUDPAddr)
	if err != nil {
		log.Error(fmt.Errorf("dialling vpn: %v", err), nil)
		os.Exit(1)
	}
	defer vpnConn.Close()

	log.Info(fmt.Sprintf("VPN UDP conn: %s", vpnAddr), nil)

	serverAddr := net.JoinHostPort(serverIP, serverPort)

	switch protocol {
	case "tcp":
		tcpConn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			log.Error(fmt.Errorf("dialling tcp server: %v", err), nil)
			os.Exit(1)
		}
		defer tcpConn.Close()

		log.Info("tcp connected.", nil)
		log.Info("proxy started...", nil)

		go io.Copy(tcpConn, clientUDPConn)
		go io.Copy(vpnConn, tcpConn)
		go io.Copy(tcpConn, vpnConn)
	default:
		wsEndpoint := fmt.Sprintf("ws://%s/ws", serverAddr)
		wsConn, _, err := websocket.Dial(context.Background(), wsEndpoint, nil)
		if err != nil {
			log.Error(fmt.Errorf("dialling ws: %v", err), nil)
			os.Exit(1)
		}
		defer wsConn.Close(websocket.StatusNormalClosure, "")

		wsNetConn := websocket.NetConn(context.Background(), wsConn, websocket.MessageBinary)

		log.Info("ws connected.", nil)
		log.Info("proxy started...", nil)

		go io.Copy(wsNetConn, clientUDPConn)
		go io.Copy(vpnConn, wsNetConn)
		go io.Copy(wsNetConn, vpnConn)
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
