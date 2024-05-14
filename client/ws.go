package main

import (
	"context"
	"fmt"
	"io"
	"net"

	"nhooyr.io/websocket"

	"tunelo/pkg/logger"
)

type wsTransport struct {
	serverAddr    string
	clientUDPConn *net.UDPConn
	vpnConn       *net.UDPConn
	logger        logger.Logger
}

func (t *wsTransport) run() error {
	wsEndpoint := fmt.Sprintf("ws://%s/ws", t.serverAddr)
	wsConn, _, err := websocket.Dial(context.Background(), wsEndpoint, nil)
	if err != nil {
		return fmt.Errorf("error dialling ws: %v", err)
	}
	defer func(wsConn *websocket.Conn, code websocket.StatusCode, reason string) {
		err := wsConn.Close(code, reason)
		if err != nil {
			t.logger.Error(fmt.Errorf("error closing ws conn: %v", err), nil)
		}
	}(wsConn, websocket.StatusNormalClosure, "")

	wsNetConn := websocket.NetConn(context.Background(), wsConn, websocket.MessageBinary)

	t.logger.Info("ws connected. Tunneling...", nil)

	go func() {
		_, err := io.Copy(wsNetConn, t.clientUDPConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from client to ws conn: %v", err), nil)
		}
	}()

	go func() {
		_, err := io.Copy(t.vpnConn, wsNetConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from ws conn to vpn: %v", err), nil)
		}
	}()

	go func() {
		_, err := io.Copy(wsNetConn, t.vpnConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying from vpn to ws conn: %v", err), nil)
		}
	}()

	return nil
}
