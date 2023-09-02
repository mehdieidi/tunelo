package ws

import (
	"context"
	"fmt"
)

func (s *WebSocket) Read() {
	for {
		_, buf, err := s.Conn.Read(context.Background())
		if err != nil {
			s.Logger.Error(fmt.Errorf("[error] reading from ws: %v", err), nil)
			break
		}

		s.Logger.Info("[info] data received over ws.", nil)

		go s.MsgHandlerFunc(buf)
	}
}
