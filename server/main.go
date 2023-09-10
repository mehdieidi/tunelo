package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
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
	flag.StringVar(&protocol, "p", "ws", "Tunnel transport protocol. Options: ws and tcp.")
	flag.Parse()

	log := plain.New()

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

	serverAddr := net.JoinHostPort(serverIP, serverPort)

	switch protocol {
	case "tcp":
		tcpListener, err := net.Listen("tcp", serverAddr)
		if err != nil {
			log.Error(fmt.Errorf("creating tcp listener: %v", err), nil)
			os.Exit(1)
		}
		defer tcpListener.Close()

		for {
			tcpConn, err := tcpListener.Accept()
			if err != nil {
				log.Error(fmt.Errorf("accepting tcp conn: %v", err), nil)
				break
			}
			defer tcpConn.Close()

			go io.Copy(vpnConn, tcpConn)
			go io.Copy(tcpConn, vpnConn)
		}
	default:
		ws := ws{vpnConn: vpnConn, log: log}

		http.HandleFunc("/ws", ws.handler)

		log.Info(fmt.Sprintf("WebSocket server listening on %s", serverAddr), nil)
		if err := http.ListenAndServe(serverAddr, nil); err != nil {
			log.Error(fmt.Errorf("listen and server: %v", err), nil)
			os.Exit(1)
		}
	}
}
