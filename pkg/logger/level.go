package logger

type Level uint8

const (
	UnsetLevel Level = iota
	InfoLevel
	NoticeLevel
	DebugLevel
	DeepDebugLevel
	WarningLevel
	ErrorLevel
	AlertLevel
	PanicLevel
	CriticalLevel
	EmergencyLevel
	FatalLevel
	SecurityLevel
	ConfidentialLevel
)

func (l Level) String() string {
	switch l {
	case UnsetLevel:
		return "UNSET"
	case InfoLevel:
		return "INFO"
	case NoticeLevel:
		return "NOTICE"
	case DebugLevel:
		return "DEBUG"
	case DeepDebugLevel:
		return "DEEP_DEBUG"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case AlertLevel:
		return "ALERT"
	case PanicLevel:
		return "PANIC"
	case CriticalLevel:
		return "CRITICAl"
	case EmergencyLevel:
		return "EMERGENCY"
	case FatalLevel:
		return "FATAL"
	case SecurityLevel:
		return "SECURITY"
	case ConfidentialLevel:
		return "CONFIDENTIAL"
	default:
		return ""
	}
}
