package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"donatello/pkg/xcrypto"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *handler) sendToWS(data []byte) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.wsConn.WriteMessage(websocket.BinaryMessage, data)
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
			break
		}

		fmt.Println("[info] read from ws.")

		go h.wsMsgHandler(buf)
	}
}

func (h *handler) wsMsgHandler(msg []byte) {
	decryptedData, err := xcrypto.Decrypt(msg, h.secretKey)
	if err != nil {
		fmt.Println(err)
		h.logFile.WriteString(err.Error() + "\n")
		return
	}

	h.sendToWireguard(decryptedData)
}
