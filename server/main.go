package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"nhooyr.io/websocket"
)

type ws struct {
	udpConn *net.UDPConn
}

func (s *ws) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println("accepting ws conn", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	wsNetConn := websocket.NetConn(context.Background(), conn, websocket.MessageBinary)

	go io.Copy(s.udpConn, wsNetConn)
	io.Copy(wsNetConn, s.udpConn)
}

func main() {
	var serverIP string
	var serverPort string
	var vpnPort string

	flag.StringVar(&serverIP, "server_ip", "127.0.0.1", "Server IP.")
	flag.StringVar(&serverPort, "server_port", "23230", "Server port.")
	flag.StringVar(&vpnPort, "vpn_port", "23233", "VPN port.")
	flag.Parse()

	serverAddr := net.JoinHostPort(serverIP, serverPort)
	vpnAddr := net.JoinHostPort("127.0.0.1", vpnPort)

	udpAddr, err := net.ResolveUDPAddr("udp", vpnAddr)
	if err != nil {
		log.Fatal("resolving udp addr", err)
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:") // Use any available local port
	if err != nil {
		log.Fatal("resolving local UDP addr:", err)
	}

	udpConn, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		log.Fatal("dialling vpn", err)
	}
	defer udpConn.Close()

	ws := ws{udpConn: udpConn}

	http.HandleFunc("/ws", ws.handler)

	fmt.Println("WebSocket server listening on", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("websocket server listening", err)
	}
}
