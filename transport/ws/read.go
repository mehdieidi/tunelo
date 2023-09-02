package ws

import (
	"context"
	"fmt"
)

func (s *WebSocket) Read() {
	if s.MsgHandlerFunc == nil {
		s.Logger.Error(fmt.Errorf("no websocket msg handler registered"), nil)
		return
	}
	if s.Conn == nil {
		s.Logger.Error(fmt.Errorf("no websocket conn registered"), nil)
		return
	}

	for {
		_, buf, err := s.Conn.Read(context.Background())
		if err != nil {
			s.Logger.Error(fmt.Errorf("reading from ws: %v", err), nil)
			break
		}

		s.Logger.Info("data received over ws.", nil)

		go s.MsgHandlerFunc(buf)
	}
}
