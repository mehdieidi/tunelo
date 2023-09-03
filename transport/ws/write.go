package ws

import (
	"context"
	"time"

	"nhooyr.io/websocket"

	"tunelo/transport"
)

func Write(conn *transport.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return conn.WebSocket.Write(ctx, websocket.MessageBinary, msg)
}
