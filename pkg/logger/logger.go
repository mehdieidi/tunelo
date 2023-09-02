package logger

const (
	LogErrKey  = "err"
	LogRespKey = "response"
	LogObjKey  = "obj"
)

const (
	DomainJSONKey = "domain"
	LayerJSONKey  = "layer"
	MethodJSONKey = "method"
	TraceJSONKey  = "trace"
	LevelJSONKey  = "level"
	FileJSONKey   = "file"
	LineJSONKey   = "line"
	CallerJSONKey = "caller"
	MsgJSONKey    = "msg"
)

type Args map[string]any

type Logger interface {
	PanicHandler()
	Info(string, Args)
	Error(error, Args)
	Panic(Args)
}
