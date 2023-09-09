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

	flag.StringVar(&serverIP, "server_ip", "127.0.0.1", "Server IP.")
	flag.StringVar(&serverPort, "server_port", "23230", "Server port.")
	flag.Parse()

	serverAddr := net.JoinHostPort(serverIP, serverPort)
	clientAddr := net.JoinHostPort("127.0.0.1", "23231")

	udpAddr, err := net.ResolveUDPAddr("udp", clientAddr)
	if err != nil {
		log.Fatal("resolving udp addr", err)
	}

	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("creating udp listener", err)
	}
	defer listener.Close()

	wsEndpoint := fmt.Sprintf("ws://%s/ws", serverAddr)
	wsConn, _, err := websocket.Dial(context.Background(), wsEndpoint, nil)
	if err != nil {
		log.Fatal("dialling ws", err)
	}
	defer wsConn.Close(websocket.StatusNormalClosure, "")

	wsNetConn := websocket.NetConn(context.Background(), wsConn, websocket.MessageBinary)

	fmt.Println("ws and udp connected, copying...")

	go io.Copy(listener, wsNetConn)
	go io.Copy(wsNetConn, listener)

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
