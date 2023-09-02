package main

import (
	"fmt"

	"tunelo/pkg/xcrypto"
)

func (h *handler) udpReadHandler() {
	buf := make([]byte, 1450)

	for {
		n, _, err := h.udpListener.ReadFrom(buf)
		if err != nil {
			errStr := fmt.Sprintf("[error] reading data from the udp conn: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			break
		}

		fmt.Println("[info] read data from udp")

		go h.udpMsgHandler(buf[:n])
	}
}

func (h *handler) udpMsgHandler(msg []byte) {
	encryptedData, err := xcrypto.Encrypt(msg, h.secretKey)
	if err != nil {
		fmt.Println(err)
		h.logFile.WriteString(err.Error() + "\n")
		return
	}

	if err = h.sendToWS(encryptedData); err != nil {
		errStr := fmt.Sprintf("[error] ws write: %v\n", err.Error())
		fmt.Println(errStr)
		h.logFile.WriteString(errStr)
		return
	}

	h.logFile.WriteString("[info] data sent to ws.\n")
}
