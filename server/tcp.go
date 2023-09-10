package main

import (
	"fmt"
	"io"
	"net"

	"tunelo/pkg/logger"
)

type tcpTransport struct {
	vpnConn *net.UDPConn
	log     logger.Logger
}

func (t *tcpTransport) handle(conn net.Conn) {
	defer conn.Close()

	go func() {
		if _, err := io.Copy(t.vpnConn, conn); err != nil {
			t.log.Error(fmt.Errorf("copying from tcp conn to vpn: %v", err), nil)
		}
	}()

	if _, err := io.Copy(conn, t.vpnConn); err != nil {
		t.log.Error(fmt.Errorf("copying from vpn to tcp conn: %v", err), nil)
	}
}
