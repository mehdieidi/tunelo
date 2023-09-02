package main

import (
	"fmt"
	"time"

	"tunelo/pkg/xcrypto"

	"github.com/gorilla/websocket"
)

func connectServerWS(serverAddr string) (*websocket.Conn, error) {
	wsConn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	return wsConn, err
}

func tryConnectServerWS(serverAddr string) (*websocket.Conn, error) {
	const tryLimit = 10

	for i := 0; i < tryLimit; i++ {
		wsConn, err := connectServerWS(serverAddr)
		if err == nil {
			return wsConn, nil
		}

		fmt.Println("[error] connecting server ws. waiting for 5 seconds")

		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("[error] connecting server ws")
}

func (h *handler) sendToWS(data []byte) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.wsConn.WriteMessage(websocket.BinaryMessage, data)
}

// wsReadHandler is supposed to handle the response data coming from the server websocket.
func (h *handler) wsReadHandler() {
	for {
		_, msg, err := h.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) {
				wsConn, err := tryConnectServerWS(h.wsServerAddr)
				if err != nil {
					fmt.Println(err)
					h.logFile.WriteString(err.Error() + "\n")
					break
				}
				h.wsConn = wsConn
				continue
			} else {
				errStr := fmt.Sprintf("[error] reading data from ws: %v\n", err.Error())
				fmt.Println(errStr)
				h.logFile.WriteString(errStr)
				break
			}
		}

		fmt.Println("[info] read from ws")

		go h.wsMsgHandler(msg)
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
