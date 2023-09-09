package wire

import (
	"fmt"

	"tunelo/pkg/xcrypto"
	"tunelo/transport"
	"tunelo/transport/ws"
)

func (w *Wire) UDPMsgHandler(transportConn *transport.Conn, msg []byte) {
	encryptedMsg, err := xcrypto.Encrypt(msg, w.SecretKey)
	if err != nil {
		w.Logger.Error(fmt.Errorf("encrypting: %v", err), nil)
		return
	}

	if err := ws.Write(transportConn, encryptedMsg); err != nil {
		w.Logger.Error(fmt.Errorf("writing to ws: %v", err), nil)
		return
	}

	w.Logger.Info("response written to ws.", nil)
}
