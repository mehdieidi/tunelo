package main

import (
	"flag"
	"fmt"
	"net"

	"golang.org/x/crypto/chacha20poly1305"
)

var secretKey = []byte("123abc@#$456asdfg#$%89756*&^fegv")

func main() {
	serverIPFlag := flag.String("i", "127.0.0.1", "Server IP address.")
	serverPortFlag := flag.String("p", "23230", "Server port.")
	flag.Parse()

	listener, err := net.Listen("tcp", *serverIPFlag+":"+*serverPortFlag)
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening on", *serverIPFlag, ":", *serverPortFlag)

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
