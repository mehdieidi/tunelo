package main

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

func encrypt(buf, secretKey []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("[error] generating nonce: " + err.Error())
	}

	aead, err := chacha20poly1305.New(secretKey)
	if err != nil {
		return nil, fmt.Errorf("[error] creating aead: " + err.Error())
	}

	encryptedData := aead.Seal(nil, nonce, buf, nil)
	encryptedDataWithNonce := append(nonce, encryptedData...)

	return encryptedDataWithNonce, nil
}
