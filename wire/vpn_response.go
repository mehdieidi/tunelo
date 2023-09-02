package wire

import (
	"fmt"
	"net"

	"tunelo/pkg/xcrypto"
)

func (w *Wire) readVPNResponse(conn net.Conn) {
	buf := make([]byte, w.BufSize)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			w.Logger.Error(fmt.Errorf("reading from vpn conn: %v", err), nil)
			break
		}

		w.Logger.Info("read vpn response.", nil)

		go w.vpnResponseMsgHandler(buf[:n])
	}

	conn.Close()

	w.Logger.Info("closed vpn udp conn.", nil)
}

func (w *Wire) vpnResponseMsgHandler(data []byte) {
	encryptedData, err := xcrypto.Encrypt(data, w.SecretKey)
	if err != nil {
		w.Logger.Error(fmt.Errorf("encrypting: %v", err), nil)
		return
	}

	if err := w.WebSocket.Write(encryptedData); err != nil {
		w.Logger.Error(fmt.Errorf("writing to ws: %v", err), nil)
		return
	}

	w.Logger.Info("response written to ws.", nil)
}
