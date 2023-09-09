package ws

import (
	"fmt"
	"net/http"

	"nhooyr.io/websocket"

	"tunelo/transport"
)

func (s *WebSocket) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.Logger.Error(fmt.Errorf("accepting websocket connection: %v", err), nil)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	c := &transport.Conn{WebSocket: conn}

	if err := s.Read(c); err != nil {
		s.Logger.Error(fmt.Errorf("reading ws conn: %v", err), nil)
		return
	}
}
