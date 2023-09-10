package main

import (
	"fmt"
	"io"
	"net"

	"tunelo/pkg/logger"
)

type utlsTransport struct {
	vpnConn *net.UDPConn
	log     logger.Logger
}

func (t *utlsTransport) handle(conn net.Conn) {
	defer conn.Close()

	go func() {
		if _, err := io.Copy(t.vpnConn, conn); err != nil {
			t.log.Error(fmt.Errorf("copying from tls conn to vpn: %v", err), nil)
		}
	}()

	if _, err := io.Copy(conn, t.vpnConn); err != nil {
		t.log.Error(fmt.Errorf("copying from vpn to tls conn: %v", err), nil)
	}
}
