package main

import (
	"fmt"

	"donatello/pkg/xcrypto"
)

func (h *handler) wsReadHandler() {
	for {
		_, msg, err := h.wsConn.ReadMessage()
		if err != nil {
			errStr := fmt.Sprintf("[error] reading data from ws: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			continue
		}

		fmt.Println("[info] read data from ws.")

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
