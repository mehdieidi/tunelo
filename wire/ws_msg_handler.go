package wire

import (
	"fmt"

	"tunelo/pkg/xcrypto"
)

func (w *Wire) WebSocketMsgHandler(msg []byte) {
	decryptedData, err := xcrypto.Decrypt(msg, w.SecretKey)
	if err != nil {
		w.Logger.Error(fmt.Errorf("decrypting: %v", err), nil)
		return
	}

	w.DialVPN(decryptedData)
}
