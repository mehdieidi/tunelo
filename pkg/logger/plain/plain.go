package plain

import (
	"fmt"
	"io"
	"runtime/debug"

	"tunelo/pkg/logger"
)

type plain struct{}

func New(w io.Writer) logger.Logger {
	return &plain{}
}

func (z plain) PanicHandler() {
	if r := recover(); r != nil {
		z.Panic(logger.Args{"err": r})
	}
}

func (z plain) Info(msg string, args logger.Args) {
	msg = fmt.Sprintf("[info] %s", msg)
	fmt.Println(msg)
}

func (z plain) Error(err error, args logger.Args) {
	msg := fmt.Sprintf("[error] %s", err.Error())
	fmt.Println(msg)
}

func (z plain) Panic(args logger.Args) {
	msg := fmt.Sprintf("[panic] %s", string(debug.Stack()))
	fmt.Println(msg)
}
