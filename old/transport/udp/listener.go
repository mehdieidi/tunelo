package udp

import (
	"net"
)

func (u *UDP) Listen() error {
	listener, err := net.ListenPacket("udp", u.ServerAddr)
	if err != nil {
		return err
	}

	u.Listener = listener

	return nil
}
