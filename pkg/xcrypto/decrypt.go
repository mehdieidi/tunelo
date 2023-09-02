package xcrypto

import (
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

func Decrypt(cipher, secret []byte) ([]byte, error) {
	nonce := cipher[:chacha20poly1305.NonceSize]
	encryptedData := cipher[chacha20poly1305.NonceSize:]

	aead, err := chacha20poly1305.New(secret)
	if err != nil {
		return nil, fmt.Errorf("error creating aead: %v", err)
	}

	decryptedData, err := aead.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting: %v", err)
	}

	return decryptedData, nil
}
