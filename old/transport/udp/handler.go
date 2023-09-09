package udp

import (
	"fmt"

	"tunelo/transport"
)

func (u *UDP) Read(transportConn *transport.Conn) {
	if u.Listener == nil {
		u.Logger.Error(fmt.Errorf("no udp listener registered"), nil)
		return
	}

	defer u.Listener.Close()

	if u.MsgHandlerFunc == nil {
		u.Logger.Error(fmt.Errorf("no udp msg handler registered"), nil)
		return
	}

	buf := make([]byte, u.BufSize)

	for {
		n, _, err := u.Listener.ReadFrom(buf)
		if err != nil {
			u.Logger.Error(fmt.Errorf("reading udp listener data: %v", err), nil)
			return
		}

		u.Logger.Info("read from udp listener.", nil)

		go u.MsgHandlerFunc(transportConn, buf[:n])
	}
}
