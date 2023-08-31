package main

import (
	"net"
	"os"

	"github.com/gorilla/websocket"
)

type handler struct {
	wgPort      string
	wsConn      *websocket.Conn
	udpListener net.PacketConn
	secretKey   []byte
	logFile     *os.File
}
