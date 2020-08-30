package gerr

// Log interface for log
type Log interface {
	Log(vals ...interface{}) error
	Debug(vals ...interface{}) error
	Info(vals ...interface{}) error
	Warn(vals ...interface{}) error
	Error(vals ...interface{}) error
	Errorf(str string, vals ...interface{}) error
}
