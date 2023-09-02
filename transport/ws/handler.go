package ws

import (
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

func (s *WebSocket) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.Logger.Error(fmt.Errorf("[error] accepting websocket connection: %v", err), nil)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "")

	s.Conn = conn

	for {
		_, buf, err := conn.Read(r.Context())
		if err != nil {
			s.Logger.Error(fmt.Errorf("[error] reading from ws: %v", err), nil)
			break
		}

		s.Logger.Info("[info] data received over ws.", nil)

		go s.MsgHandlerFunc(buf)
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
