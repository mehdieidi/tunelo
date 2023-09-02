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
			w.Logger.Error(fmt.Errorf("[error] reading from vpn conn: %v", err), nil)
			break
		}

		w.Logger.Info("[info] read vpn response.", nil)

		go w.vpnResponseMsgHandler(buf[:n])
	}

	conn.Close()

	w.Logger.Info("[info] closed vpn udp conn.", nil)
}

func (w *Wire) vpnResponseMsgHandler(data []byte) {
	encryptedData, err := xcrypto.Encrypt(data, w.SecretKey)
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
