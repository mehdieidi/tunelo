package wire

import (
	"fmt"
	"net"
)

func (w *Wire) DialVPN(data []byte) {
	conn, err := net.Dial("udp", w.VPNAddr)
	if err != nil {
		w.Logger.Error(fmt.Errorf("[error] dialing vpn udp conn: %v", err), nil)
		return
	}

	w.Logger.Info("[info] dialed vpn.", nil)

	go w.readVPNResponse(conn)

	_, err = conn.Write(data)
	if err != nil {
		w.Logger.Error(fmt.Errorf("[error] writing to vpn conn: %v", err), nil)
		return
	}

	w.Logger.Info("[info] sent data to vpn.", nil)
}
