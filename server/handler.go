package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsHandler struct {
	secretKey []byte
	logFile   *os.File
}

func (wh *wsHandler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errStr := "[error] ws upgrade: " + err.Error()
		fmt.Println(errStr)
		wh.logFile.WriteString(errStr + "\n")
		return
	}
	defer conn.Close()

	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			errStr := "[error] ws read: " + err.Error()
			fmt.Println(errStr)
			wh.logFile.WriteString(errStr + "\n")
			continue
		}

		nonce := buf[:chacha20poly1305.NonceSize]
		encryptedData := buf[chacha20poly1305.NonceSize:]

		aead, err := chacha20poly1305.New(wh.secretKey)
		if err != nil {
			errStr := "[error] creating aead: " + err.Error()
			fmt.Println(errStr)
			wh.logFile.WriteString(errStr + "\n")
			continue
		}

		decryptedData, err := aead.Open(nil, nonce, encryptedData, nil)
		if err != nil {
			errStr := "[error] decrypting: " + err.Error()
			fmt.Println(errStr)
			wh.logFile.WriteString(errStr + "\n")
			continue
		}

		fmt.Println("decrypted data:", string(decryptedData))
		sendToWireguard(decryptedData)
	}
}
