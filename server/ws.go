package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"nhooyr.io/websocket"

	"tunelo/pkg/logger"
)

type wsTransport struct {
	serverAddr string
	vpnConn    *net.UDPConn
	logger     logger.Logger
}

func (t *wsTransport) run() error {
	http.HandleFunc("/ws", t.handler)

	t.logger.Info(fmt.Sprintf("WebSocket server listening on %s", t.serverAddr), nil)

	err := http.ListenAndServe(t.serverAddr, nil)
	if err != nil {
		return fmt.Errorf("error listening: %v", err)
	}

	return nil
}

func (t *wsTransport) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		t.logger.Error(fmt.Errorf("error accepting ws conn: %v", err), nil)
		return
	}
	defer func(conn *websocket.Conn, code websocket.StatusCode, reason string) {
		err := conn.Close(code, reason)
		if err != nil {
			t.logger.Error(fmt.Errorf("error closing conn: %v", err), nil)
		}
	}(conn, websocket.StatusNormalClosure, "")

	wsNetConn := websocket.NetConn(r.Context(), conn, websocket.MessageBinary)

	t.logger.Info("ws connection accepted. Tunneling...", nil)

	go func() {
		_, err := io.Copy(t.vpnConn, wsNetConn)
		if err != nil {
			t.logger.Error(fmt.Errorf("error copying data: %v", err), nil)
		}
	}()

	_, err = io.Copy(wsNetConn, t.vpnConn)
	if err != nil {
		t.logger.Error(fmt.Errorf("error copying data: %v", err), nil)
	}
}
