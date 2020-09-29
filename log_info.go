package gerr

import "fmt"

const (
	// LogKeyService log key service
	LogKeyService = "service"

	// LogKeyEnvironment log key Environment
	LogKeyEnvironment = "env"

	// LogKeyTraceID log key TraceId
	LogKeyTraceID = "traceId"

	// LogKeyPath log key Path
	LogKeyPath = "path"

	// LogKeyMethod log key Method
	LogKeyMethod = "method"

	// LogKeyIP log key User IP
	LogKeyIP = "ip"

	// LogKeyUserAgent log key user agent
	LogKeyUserAgent = "userAgent"

	// LogKeyUnknown log key Method
	LogKeyUnknown = "unknown"
)

// Service service name for logger
type Service string

// Environment service name for logger
type Environment string

// RequestInfo request information
type RequestInfo struct {
	TraceID   string
	Path      string
	Method    string
	IP        string
	UserAgent string
}

// LogInfo base info for log
type LogInfo map[string]string

// Set set value for map
func (l LogInfo) Set(key, val string) {
	l[key] = val
}

// Update update value for map
func (l LogInfo) Update(args ...interface{}) {
	for _, arg := range args {
		switch arg := arg.(type) {
		case Service:
			l[LogKeyService] = string(arg)
			break
		case Environment:
			l[LogKeyEnvironment] = string(arg)
		case RequestInfo:
			l[LogKeyTraceID] = arg.TraceID
			l[LogKeyPath] = arg.Path
			l[LogKeyMethod] = arg.Method
			l[LogKeyUserAgent] = arg.UserAgent
			l[LogKeyIP] = arg.IP

		default:
			l[LogKeyUnknown] = fmt.Sprintf("%v", arg)
		}
	}
}

// GetTraceID get trace ID
func (l LogInfo) GetTraceID() string {
	if val, ok := l[LogKeyTraceID]; ok {
		return val
	}

	return ""
}

// NewLogInfo make base information for log
func NewLogInfo(args ...interface{}) LogInfo {
	rs := LogInfo{}
	rs.Update(args...)
	return rs
}
