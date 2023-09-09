package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"nhooyr.io/websocket"
)

func main() {
	var serverIP string
	var serverPort string
	var vpnPort string

	flag.StringVar(&serverIP, "server_ip", "127.0.0.1", "Server IP.")
	flag.StringVar(&serverPort, "server_port", "23230", "Server port.")
	flag.StringVar(&vpnPort, "vpn_port", "23233", "VPN port.")
	flag.Parse()

	clientAddr := net.JoinHostPort("127.0.0.1", "23231")
	clientUDPAddr, err := net.ResolveUDPAddr("udp", clientAddr)
	if err != nil {
		log.Fatal("resolving udp addr ", err)
	}

	clientUDPConn, err := net.ListenUDP("udp", clientUDPAddr)
	if err != nil {
		log.Fatal("creating udp listener ", err)
	}
	defer clientUDPConn.Close()

	vpnAddr := net.JoinHostPort("127.0.0.1", vpnPort)
	vpnUDPAddr, err := net.ResolveUDPAddr("udp", vpnAddr)
	if err != nil {
		log.Fatal("resolving udp addr ", err)
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:")
	if err != nil {
		log.Fatal("resolving local UDP addr:", err)
	}

	vpnConn, err := net.DialUDP("udp", localAddr, vpnUDPAddr)
	if err != nil {
		log.Fatal("creating udp listener ", err)
	}
	defer vpnConn.Close()

	serverAddr := net.JoinHostPort(serverIP, serverPort)
	wsEndpoint := fmt.Sprintf("ws://%s/ws", serverAddr)
	wsConn, _, err := websocket.Dial(context.Background(), wsEndpoint, nil)
	if err != nil {
		log.Fatal("dialling ws ", err)
	}
	defer wsConn.Close(websocket.StatusNormalClosure, "")

	wsNetConn := websocket.NetConn(context.Background(), wsConn, websocket.MessageBinary)

	fmt.Println("ws and udp connected, copying...")

	go io.Copy(wsNetConn, clientUDPConn)
	go io.Copy(vpnConn, wsNetConn)

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
