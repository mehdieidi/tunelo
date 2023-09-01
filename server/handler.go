package main

import (
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type handler struct {
	wgAddr    string
	secretKey []byte
	logFile   *os.File
	wsConn    *websocket.Conn
	mutex     sync.Mutex
}
