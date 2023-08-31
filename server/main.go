package main

import (
	"fmt"
	"net"

	"golang.org/x/crypto/chacha20poly1305"
)

var secretKey = []byte("123abc@#$456asdfg#$%89756*&^fegv")

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:23230")
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 2048)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error connecting:", err)
		}

		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("error reading from connection:", err)
				break
			}

			fmt.Println("received from client")

			go handle(buf[:n])
		}
	}
}

func handle(buf []byte) {
	nonce := buf[:chacha20poly1305.NonceSize]
	encryptedData := buf[chacha20poly1305.NonceSize:]

	aead, err := chacha20poly1305.New(secretKey)
	if err != nil {
		fmt.Println("error creating aead:", err)
		return
	}

	decryptedData, err := aead.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		fmt.Println("error decrypting:", err)
		return
	}

	fmt.Println("decrypted data:", string(decryptedData))
}
