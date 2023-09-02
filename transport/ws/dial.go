package ws

import (
	"context"
	"time"

	"nhooyr.io/websocket"
)

func Dial(addr string) (*websocket.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, addr, nil)

	return conn, err
}
