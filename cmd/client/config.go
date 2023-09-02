package client

import "tunelo/pkg/logger"

type Config struct {
	VPNPort          string
	ClientPort       string
	RemoteServerIP   string
	RemoteServerPort string
	SecretKey        []byte
	Logger           logger.Logger
}
