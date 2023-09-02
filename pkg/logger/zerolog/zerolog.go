package zerolog

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/rs/zerolog"

	"tunelo/pkg/logger"
)

type zeroLog struct {
	logger zerolog.Logger
}

func New(w io.Writer) logger.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zeroLog{logger: zerolog.New(w).With().Timestamp().Logger()}
}

func (z zeroLog) PanicHandler() {
	if r := recover(); r != nil {
		z.Panic(logger.Args{"err": r})
	}
}

func (z zeroLog) Info(msg string, args logger.Args) {
	funcName, _, _ := logger.Caller()

	e := z.logger.Info().
		Str(logger.CallerJSONKey, funcName).
		Str(logger.MsgJSONKey, msg)

	for k, v := range args {
		if k == logger.LogObjKey {
			j, _ := json.Marshal(v)
			e.RawJSON(logger.LogObjKey, j)
		} else {
			e.Str(k, fmt.Sprintf("%+v", v))
		}
	}

	e.Msg("")
}

func (z zeroLog) Error(err error, args logger.Args) {
	funcName, file, line := logger.Caller()

	e := z.logger.Error().
		Str(logger.FileJSONKey, file).
		Int(logger.LineJSONKey, line).
		Str(logger.CallerJSONKey, funcName).
		Err(err)

	for k, v := range args {
		e.Str(k, fmt.Sprintf("%+v", v))
	}

	e.Msg("")
}

func (z zeroLog) Panic(args logger.Args) {
	funcName, file, line := logger.Caller()

	e := z.logger.Log().
		Str(logger.LevelJSONKey, "panic").
		Str(logger.TraceJSONKey, string(debug.Stack())).
		Str(logger.FileJSONKey, file).
		Int(logger.LineJSONKey, line).
		Str(logger.CallerJSONKey, funcName)

	for k, v := range args {
		e.Str(k, fmt.Sprintf("%+v", v))
	}

	e.Msg("")
}
