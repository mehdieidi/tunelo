package main

import (
	"fmt"
	"net"
)

func sendToWireguard(data []byte) {
	conn, err := net.Dial("udp", "127.0.0.1:51820")
	if err != nil {
		fmt.Println("Error connecting to WireGuard:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}

	fmt.Println("data sent to wg")
}
