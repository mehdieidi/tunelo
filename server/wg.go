package main

import (
	"fmt"
	"net"

	"github.com/gorilla/websocket"

	"donatello/pkg/xcrypto"
)

func (h *handler) sendToWireguard(data []byte) {
	wgAddr := "127.0.0.1" + ":" + h.wgPort
	conn, err := net.Dial("udp", wgAddr)
	if err != nil {
		errStr := fmt.Sprintf("[error] dialing wg: %v\n", err.Error())
		fmt.Println(errStr)
		h.logFile.WriteString(errStr)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		errStr := fmt.Sprintf("[error] writing to wg conn: %v\n", err.Error())
		fmt.Println(errStr)
		h.logFile.WriteString(errStr)
		return
	}

	fmt.Println("[info] data sent to wg.")

	h.readFromWireguard(conn)
}

func (h *handler) readFromWireguard(conn net.Conn) {
	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			errStr := fmt.Sprintf("[error] reading from wg conn: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			break
		}

		encryptedData, err := xcrypto.Encrypt(buf[:n], h.secretKey)
		if err != nil {
			errStr := fmt.Sprintf("[error] encrypting: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			break
		}

		if err := h.wsConn.WriteMessage(websocket.BinaryMessage, encryptedData); err != nil {
			errStr := fmt.Sprintf("[error] writing data to ws: %v\n", err.Error())
			fmt.Println(errStr)
			h.logFile.WriteString(errStr)
			break
		}
	}
}
