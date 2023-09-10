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
	vpnConn *net.UDPConn
	log     logger.Logger
}

func (s *wsTransport) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.log.Error(fmt.Errorf("accepting ws conn: %v", err), nil)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	wsNetConn := websocket.NetConn(r.Context(), conn, websocket.MessageBinary)

	s.log.Info("ws connection accepted. Proxy started...", nil)

	go io.Copy(s.vpnConn, wsNetConn)
	io.Copy(wsNetConn, s.vpnConn)
}
