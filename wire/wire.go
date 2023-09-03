package wire

import (
	"tunelo/pkg/logger"
)

type Wire struct {
	SecretKey []byte
	Logger    logger.Logger
	VPNAddr   string
	BufSize   int
}
