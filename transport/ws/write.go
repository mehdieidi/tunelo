package ws

import (
	"context"
	"time"

	"nhooyr.io/websocket"
)

func (s *WebSocket) Write(msg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Conn.Write(ctx, websocket.MessageBinary, msg)
}
