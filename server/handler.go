package main

import (
	"donatello/pkg/xcrypto"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type handler struct {
	wgPort    string
	secretKey []byte
	logFile   *os.File
	wsConn    *websocket.Conn
}

func (h *handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		errStr := fmt.Sprintf("[error] ws upgrade: %v\n", err.Error())
		fmt.Println(errStr)
		h.logFile.WriteString(errStr)
		return
	}
	defer conn.Close()

	h.wsConn = conn

	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			errStr := fmt.Sprintf("[error] ws read: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			continue
		}

		decryptedData, err := xcrypto.Decrypt(buf, h.secretKey)
		if err != nil {
			fmt.Println(err)
			h.logFile.WriteString(err.Error() + "\n")
			continue
		}

		fmt.Println("decrypted data from ws:", string(decryptedData))

		h.sendToWireguard(decryptedData)
	}
}
