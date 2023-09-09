package ws

import (
	"fmt"
	"net/http"
)

func (s *WebSocket) ListenAndServe() error {
	http.HandleFunc("/"+s.Endpoint, s.handler)

	s.Logger.Info(fmt.Sprintf("WebSocket server listening on %s", s.ServerAddr), nil)

	return http.ListenAndServe(s.ServerAddr, nil)
}
