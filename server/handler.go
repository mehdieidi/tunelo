package main

import (
	"os"

	"github.com/gorilla/websocket"
)

type handler struct {
	wgPort    string
	secretKey []byte
	logFile   *os.File
	wsConn    *websocket.Conn
}
