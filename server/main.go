package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"nhooyr.io/websocket"

	"tunelo/pkg/logger"
	"tunelo/pkg/logger/plain"
)

type ws struct {
	vpnConn *net.UDPConn
	logger  logger.Logger
}

func (s *ws) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.logger.Error(fmt.Errorf("accepting ws conn: %v", err), nil)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	wsNetConn := websocket.NetConn(r.Context(), conn, websocket.MessageBinary)

	s.logger.Info("ws connection accepted. Proxy started.", nil)

	go io.Copy(s.vpnConn, wsNetConn)
	io.Copy(wsNetConn, s.vpnConn)
}

func main() {
	var serverIP string
	var serverPort string
	var vpnPort string

	flag.StringVar(&serverIP, "server_ip", "127.0.0.1", "Proxy server IP address.")
	flag.StringVar(&serverPort, "server_port", "23230", "Proxy server port number.")
	flag.StringVar(&vpnPort, "vpn_port", "23233", "Local VPN port number.")
	flag.Parse()

	logger := plain.New()

	vpnAddr := net.JoinHostPort("127.0.0.1", vpnPort)
	vpnUDPAddr, err := net.ResolveUDPAddr("udp", vpnAddr)
	if err != nil {
		logger.Error(fmt.Errorf("resolving vpn udp addr: %v", err), nil)
		os.Exit(1)
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:")
	if err != nil {
		logger.Error(fmt.Errorf("resolving local udp addr: %v", err), nil)
		os.Exit(1)
	}

	vpnConn, err := net.DialUDP("udp", localAddr, vpnUDPAddr)
	if err != nil {
		logger.Error(fmt.Errorf("dialling vpn: %v", err), nil)
		os.Exit(1)
	}
	defer vpnConn.Close()

	ws := ws{vpnConn: vpnConn}

	http.HandleFunc("/ws", ws.handler)

	serverAddr := net.JoinHostPort(serverIP, serverPort)
	logger.Info(fmt.Sprintf("WebSocket server listening on %s", serverAddr), nil)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		logger.Error(fmt.Errorf("listen and server: %v", err), nil)
		os.Exit(1)
	}
}
