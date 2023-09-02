package wire

import (
	"fmt"

	"tunelo/pkg/xcrypto"
)

func (w *Wire) UDPMsgHandler(msg []byte) {
	encryptedData, err := xcrypto.Encrypt(msg, w.SecretKey)
	if err != nil {
		w.Logger.Error(fmt.Errorf("[error] encrypting: %v", err), nil)
		return
	}

	if err := w.WebSocket.Write(encryptedData); err != nil {
		w.Logger.Error(fmt.Errorf("[error] writing to ws: %v", err), nil)
		return
	}

	w.Logger.Info("[info] response written to ws.", nil)
}
