package main

import (
	"fmt"
	"io"
	"net"

	"tunelo/pkg/logger"
)

type tcpTransport struct {
	serverAddr string
	vpnConn    *net.UDPConn
	logger     logger.Logger
}

func (t *tcpTransport) run() error {
	tcpListener, err := net.Listen("tcp", t.serverAddr)
	if err != nil {
		return fmt.Errorf("error creating tcp listener: %v", err)
	}
	defer func(c net.Listener) {
		err := c.Close()
		if err != nil {
			t.logger.Error(fmt.Errorf("error closing tcp conn: %v", err), nil)
		}
	}(tcpListener)

	t.logger.Info(fmt.Sprintf("TCP server listening on %s", t.serverAddr), nil)

	for {
		tcpConn, err := tcpListener.Accept()
		if err != nil {
			t.logger.Error(fmt.Errorf("error accepting tcp conn: %v", err), nil)
			continue
		}

		t.logger.Info("tcp connection accepted. Tunneling...", nil)

		go t.handle(tcpConn)
	}
}

func (t *tcpTransport) handle(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			t.logger.Error(fmt.Errorf("error closing tcp conn: %v", err), nil)
		}
	}(conn)

	go func() {
		_, err := io.Copy(t.vpnConn, conn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from tcp conn to vpn: %v", err), nil)
		}
	}()

	_, err := io.Copy(conn, t.vpnConn)
	if err != nil {
		t.logger.Error(fmt.Errorf("error copying from vpn to tcp conn: %v", err), nil)
	}
}
