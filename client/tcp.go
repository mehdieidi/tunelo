package main

import (
	"fmt"
	"io"
	"net"

	"tunelo/pkg/logger"
)

type tcpTransport struct {
	serverAddr    string
	vpnConn       *net.UDPConn
	clientUDPConn *net.UDPConn
	logger        logger.Logger
}

func (t *tcpTransport) run() error {
	tcpConn, err := net.Dial("tcp", t.serverAddr)
	if err != nil {
		return fmt.Errorf("error dialling tcp server: %v", err)
	}
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			t.logger.Error(fmt.Errorf("error closing tcp conn: %v", err), nil)
		}
	}(tcpConn)

	t.logger.Info("tcp connected. Tunneling...", nil)

	go func() {
		_, err := io.Copy(tcpConn, t.clientUDPConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from client to tcp conn: %v", err), nil)
		}
	}()

	go func() {
		_, err := io.Copy(t.vpnConn, tcpConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from tcp conn to vpn: %v", err), nil)
		}
	}()

	go func() {
		_, err := io.Copy(tcpConn, t.vpnConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from vpn to tcp conn: %v", err), nil)
		}
	}()

	return nil
}
