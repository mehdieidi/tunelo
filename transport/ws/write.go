package ws

import (
	"context"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

func (s *WebSocket) Write(msg []byte) error {
	if s.Conn == nil {
		return fmt.Errorf("no websocket conn registered")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return s.Conn.Write(ctx, websocket.MessageBinary, msg)
}
