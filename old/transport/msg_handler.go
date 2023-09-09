package transport

type MsgHandlerFunc func(*Conn, []byte)
