package main

import (
	"net"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type handler struct {
	wsServerAddr string
	wgAddr       string
	wsConn       *websocket.Conn
	udpListener  net.PacketConn
	secretKey    []byte
	logFile      *os.File
	mutex        sync.Mutex
}
