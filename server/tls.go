package main

import (
	"fmt"
	"io"
	"net"

	"tunelo/pkg/logger"
)

type TLS struct {
	vpnConn *net.UDPConn
	log     logger.Logger
}

func (t *TLS) handle(conn net.Conn) {
	// TODO: figure out how to handle connection closure.
	// defer conn.Close()

	go func() {
		if _, err := io.Copy(t.vpnConn, conn); err != nil {
			t.log.Error(fmt.Errorf("copying from tls conn to vpn: %v", err), nil)
		}
	}()

	go func() {
		if _, err := io.Copy(conn, t.vpnConn); err != nil {
			t.log.Error(fmt.Errorf("copying from vpn to tls conn: %v", err), nil)
		}
	}()
}
