package udp

import "fmt"

func (u *UDP) UDPReadHandler() {
	buf := make([]byte, u.BufSize)

	for {
		n, _, err := u.Conn.ReadFrom(buf)
		if err != nil {
			u.Logger.Error(fmt.Errorf("reading udp listener data: %v", err), nil)
			break
		}

		u.Logger.Info("read from udp listener.", nil)

		go u.MsgHandlerFunc(buf[:n])
	}

	u.Conn.Close()

	u.Logger.Info("closed udp listener.", nil)
}
