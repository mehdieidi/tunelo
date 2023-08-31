package main

import (
	"fmt"

	"github.com/gorilla/websocket"

	"donatello/pkg/xcrypto"
)

func (h *handler) udpReadHandler() {
	buf := make([]byte, 4096)

	for {
		n, _, err := h.udpListener.ReadFrom(buf)
		if err != nil {
			errStr := fmt.Sprintf("[error] reading data from the udp conn: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			continue
		}

		fmt.Println("received from udp conn:", string(buf))

		encryptedData, err := xcrypto.Encrypt(buf[:n], h.secretKey)
		if err != nil {
			fmt.Println(err)
			h.logFile.WriteString(err.Error() + "\n")
			continue
		}

		err = h.wsConn.WriteMessage(websocket.BinaryMessage, encryptedData)
		if err != nil {
			errStr := fmt.Sprintf("[error] ws write: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			continue
		}

		h.logFile.WriteString("[info] data sent on ws.\n")
	}
}
