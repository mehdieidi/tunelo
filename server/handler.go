package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/chacha20poly1305"
)

func handle(buf []byte, secretKey []byte, logFile *os.File) {
	nonce := buf[:chacha20poly1305.NonceSize]
	encryptedData := buf[chacha20poly1305.NonceSize:]

	aead, err := chacha20poly1305.New(secretKey)
	if err != nil {
		errStr := "[error] creating aead: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		return
	}

	decryptedData, err := aead.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		errStr := "[error] decrypting: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		return
	}

	fmt.Println("decrypted data:", string(decryptedData))
}
