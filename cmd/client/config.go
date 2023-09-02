package client

import "tunelo/pkg/logger"

type Config struct {
	VPNPort    string
	ClientPort string
	ServerIP   string
	ServerPort string
	SecretKey  []byte
	Logger     logger.Logger
	BufSize    int
}
