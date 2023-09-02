package ws

import "net/http"

func (s *WebSocket) ListenAndServe() error {
	http.HandleFunc("/ws", s.handler)
	return http.ListenAndServe(s.ServerAddr, nil)
}
