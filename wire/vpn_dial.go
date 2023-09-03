package wire

import (
	"fmt"
	"net"

	"tunelo/transport"
)

func (w *Wire) DialVPN(transportConn *transport.Conn, msg []byte) {
	vpnConn, err := net.Dial("udp", w.VPNAddr)
	if err != nil {
		w.Logger.Error(fmt.Errorf("dialing vpn udp conn: %v", err), nil)
		return
	}

	w.Logger.Info("dialed vpn.", nil)

	go w.readVPNResponse(vpnConn, transportConn)

	_, err = vpnConn.Write(msg)
	if err != nil {
		w.Logger.Error(fmt.Errorf("writing to vpn conn: %v", err), nil)
		return
	}

	w.Logger.Info("sent msg to vpn.", nil)
}
