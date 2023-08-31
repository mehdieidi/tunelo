package xcrypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(plain, secret []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("error generating nonce: " + err.Error())
	}

	aead, err := chacha20poly1305.New(secret)
	if err != nil {
		return nil, fmt.Errorf("error creating aead: " + err.Error())
	}

	encryptedData := aead.Seal(nil, nonce, plain, nil)
	encryptedDataWithNonce := append(nonce, encryptedData...)

	return encryptedDataWithNonce, nil
}
