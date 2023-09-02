package wire

import (
	"fmt"
	"net"
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

		go w.UDPMsgHandler(buf[:n])
	}

	conn.Close()

	w.Logger.Info("closed vpn udp conn.", nil)
}
