package udp

import "fmt"

func (u *UDP) Read() {
	if u.MsgHandlerFunc == nil {
		u.Logger.Error(fmt.Errorf("no udp msg handler registered"), nil)
		return
	}
	if u.Listener == nil {
		u.Logger.Error(fmt.Errorf("no udp listener registered"), nil)
		return
	}

	buf := make([]byte, u.BufSize)

	for {
		n, _, err := u.Listener.ReadFrom(buf)
		if err != nil {
			u.Logger.Error(fmt.Errorf("reading udp listener data: %v", err), nil)
			break
		}

		u.Logger.Info("read from udp listener.", nil)

		go u.MsgHandlerFunc(buf[:n])
	}

	u.Listener.Close()

	u.Logger.Info("udp listener closed.", nil)
}
