package logger

import "runtime"

func Caller() (name, file string, line int) {
	pc, file, line, ok := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	if ok && fn != nil {
		name = fn.Name()
	}
	return
}
