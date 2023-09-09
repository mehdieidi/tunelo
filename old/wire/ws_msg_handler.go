package wire

import (
	"fmt"

	"tunelo/pkg/xcrypto"
	"tunelo/transport"
)

func (w *Wire) WebSocketMsgHandler(conn *transport.Conn, msg []byte) {
	decryptedMsg, err := xcrypto.Decrypt(msg, w.SecretKey)
	if err != nil {
		w.Logger.Error(fmt.Errorf("decrypting: %v", err), nil)
		return
	}

	w.DialVPN(conn, decryptedMsg)
}
