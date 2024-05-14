package main

import (
	"fmt"
	"io"
	"net"

	tls "github.com/refraction-networking/utls"

	"tunelo/pkg/logger"
)

type utlsTransport struct {
	serverDomain  string
	serverAddr    string
	vpnConn       *net.UDPConn
	clientUDPConn *net.UDPConn
	logger        logger.Logger
}

func (t *utlsTransport) run() error {
	if t.serverDomain == "" {
		return fmt.Errorf("server domain cannot be empty")
	}

	tlsConfig := &tls.Config{
		ServerName:         t.serverDomain,
		InsecureSkipVerify: true,
	}

	tcpConn, err := net.Dial("tcp", t.serverAddr)
	if err != nil {
		return fmt.Errorf("error dialling tls server: %v", err)

	}
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			t.logger.Error(fmt.Errorf("error closing tcp conn: %v", err), nil)
		}
	}(tcpConn)

	tlsConn := tls.UClient(tcpConn, tlsConfig, tls.HelloChrome_102)
	err = tlsConn.Handshake()
	if err != nil {
		return fmt.Errorf("error in tls handshake: %v", err)
	}

	t.logger.Info("tls connected. Tunneling...", nil)

	go func() {
		_, err := io.Copy(tlsConn, t.clientUDPConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from client to tls conn: %v", err), nil)
		}
	}()

	go func() {
		_, err := io.Copy(t.vpnConn, tlsConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from tls conn to vpn: %v", err), nil)
		}
	}()

	go func() {
		_, err := io.Copy(tlsConn, t.vpnConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from vpn to tls conn: %v", err), nil)
		}
	}()

	return nil
}
