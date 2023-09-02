package udp

import (
	"fmt"
	"net"
)

func (u *UDP) Listen() error {
	conn, err := net.ListenPacket("udp", u.ServerAddr)
	if err != nil {
		u.Logger.Error(fmt.Errorf("creating udp listener: %v", err), nil)
		return err
	}

	u.Conn = conn

	return nil
}
