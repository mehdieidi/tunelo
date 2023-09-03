package transport

import "nhooyr.io/websocket"

type Conn struct {
	WebSocket *websocket.Conn
}
