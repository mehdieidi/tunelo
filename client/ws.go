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

		decryptedData, err := xcrypto.Decrypt(msg, h.secretKey)
		if err != nil {
			fmt.Println(err)
			h.logFile.WriteString(err.Error() + "\n")
			continue
		}

		fmt.Println("msg from ws:", decryptedData)

		h.sendToWireguard(decryptedData)
	}
}
