package udp

import "fmt"

func (u *UDP) UDPReadHandler() {
	buf := make([]byte, 1450)

	for {
		n, _, err := u.Conn.ReadFrom(buf)
		if err != nil {
			u.Logger.Error(fmt.Errorf("[error] reading udp listener data: %v", err), nil)
			break
		}

		u.Logger.Info("[info] read from udp listener.", nil)

		go u.MsgHandlerFunc(buf[:n])
	}

	u.Conn.Close()

	u.Logger.Info("[info] closed udp listener.", nil)
}
