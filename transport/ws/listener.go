package ws

import "net/http"

func (s *WebSocket) ListenAndServe() error {
	http.HandleFunc("/ws", s.handler)

	s.Logger.Info("[info] WebSocket server listening on "+s.ServerAddr, nil)

	return http.ListenAndServe(s.ServerAddr, nil)
}
