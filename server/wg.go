package main

import (
	"fmt"
	"net"

	"donatello/pkg/xcrypto"
)

func (h *handler) sendToWireguard(data []byte) {
	conn, err := net.Dial("udp", h.wgAddr)
	if err != nil {
		errStr := fmt.Sprintf("[error] dialing wg: %v\n", err.Error())
		fmt.Println(errStr)
		h.logFile.WriteString(errStr)
		return
	}

	fmt.Println("[info] dialed wg.")

	go h.readFromWireguard(conn)

	_, err = conn.Write(data)
	if err != nil {
		errStr := fmt.Sprintf("[error] writing to wg conn: %v\n", err.Error())
		fmt.Println(errStr)
		h.logFile.WriteString(errStr)
		return
	}

	fmt.Println("[info] data sent to wg.")
}

func (h *handler) readFromWireguard(conn net.Conn) {
	buf := make([]byte, 1380)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			errStr := fmt.Sprintf("[error] reading from wg conn: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			break
		}

		fmt.Println("[info] read data response from wg.")

		go h.wgResponseHandler(buf[:n])
	}

	conn.Close()
}

func (h *handler) wgResponseHandler(msg []byte) {
	encryptedData, err := xcrypto.Encrypt(msg, h.secretKey)
	if err != nil {
		errStr := fmt.Sprintf("[error] encrypting: %v\n", err.Error())
		fmt.Println(errStr)
		h.logFile.WriteString(errStr)
		return
	}

	if err := h.sendToWS(encryptedData); err != nil {
		errStr := fmt.Sprintf("[error] writing data to ws: %v\n", err.Error())
		fmt.Println(errStr)
		h.logFile.WriteString(errStr)
		return
	}

	fmt.Println("[info] wg response sent on ws.")
}
