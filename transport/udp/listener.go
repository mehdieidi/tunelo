package udp

import (
	"fmt"
	"net"
)

func (u *UDP) Listen() error {
	listener, err := net.ListenPacket("udp", u.ServerAddr)
	if err != nil {
		return fmt.Errorf("creating udp listener: %v", err)
	}

	u.Listener = listener

	return nil
}
