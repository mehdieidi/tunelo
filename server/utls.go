package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"tunelo/pkg/logger"
)

type utlsTransport struct {
	serverAddr string
	vpnConn    *net.UDPConn
	logger     logger.Logger
}

func (t *utlsTransport) run() error {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		return fmt.Errorf("error loading cert and key: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	tlsListener, err := tls.Listen("tcp", t.serverAddr, tlsConfig)
	if err != nil {
		return fmt.Errorf("error creating tls listener: %v", err)
	}
	defer tlsListener.Close()

	t.logger.Info(fmt.Sprintf("[+] TLS server listening on %s", t.serverAddr), nil)

	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			t.logger.Error(fmt.Errorf("error accepting tls conn: %v", err), nil)
			continue
		}

		t.logger.Info("tls connection accepted. Tunneling...", nil)

		go t.handle(conn)
	}
}

func (t *utlsTransport) handle(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			t.logger.Error(fmt.Errorf("error closing conn: %v", err), nil)
		}
	}(conn)

	go func() {
		_, err := io.Copy(t.vpnConn, conn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from tls conn to vpn: %v", err), nil)
		}
	}()

	_, err := io.Copy(conn, t.vpnConn)
	if err != nil {
		t.logger.Error(fmt.Errorf("error copying from vpn to tls conn: %v", err), nil)
	}
}
