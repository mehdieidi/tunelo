package ws

import (
	"context"
	"fmt"

	"tunelo/transport"
)

func (s *WebSocket) Read(transportConn *transport.Conn) error {
	if s.MsgHandlerFunc == nil {
		err := fmt.Errorf("no websocket msg handler registered")
		s.Logger.Error(err, nil)
		return err
	}

	for {
		_, msg, err := transportConn.WebSocket.Read(context.Background())
		if err != nil {
			return err
		}

		s.Logger.Info("msg received over ws.", nil)

		go s.MsgHandlerFunc(transportConn, msg)
	}
}
