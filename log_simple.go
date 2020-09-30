package gerr

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type logSimple struct {
	log *logrus.Logger
}

// NewSimpleLog new simple log
func NewSimpleLog() Log {
	l := logSimple{}

	var baseLogger = logrus.New()
	baseLogger.Formatter = &lokiFormatter{
		ForceQuote: true,
	}
	l.log = baseLogger

	return &l
}

func (l *logSimple) Log(vals ...interface{}) error {
	fields, lvl, vals := detachFields(vals...)
	l.log.WithFields(fields).Log(lvl, vals...)
	return nil
}
func (l *logSimple) Debug(vals ...interface{}) error {
	fields, _, vals := detachFields(vals...)
	l.log.WithFields(fields).Debug(vals...)
	return nil
}
func (l *logSimple) Info(vals ...interface{}) error {
	fields, _, vals := detachFields(vals...)
	l.log.WithFields(fields).Info(vals...)
	return nil
}
func (l *logSimple) Warn(vals ...interface{}) error {
	fields, _, vals := detachFields(vals...)
	l.log.WithFields(fields).Warning(vals...)
	return nil
}
func (l *logSimple) Error(vals ...interface{}) error {
	fields, _, vals := detachFields(vals...)
	l.log.WithFields(fields).Warning(vals...)
	return nil
}

func getLogLevel(code int) logrus.Level {
	if code < http.StatusInternalServerError {
		return logrus.InfoLevel
	}

	if code < internalCodeMax {
		return logrus.ErrorLevel
	}

	if code < serviceCodeMax {
		return logrus.ErrorLevel
	}

	if code < businessCodeMax {
		return logrus.InfoLevel
	}

	return logrus.InfoLevel
}

func detachFields(vals ...interface{}) (logrus.Fields, logrus.Level, []interface{}) {
	fields := logrus.Fields{}
	lvl := logrus.InfoLevel
	others := []interface{}{}
	for idx := range vals {
		arg := vals[idx]
		switch arg := arg.(type) {
		case logrus.Fields:
			fields = arg
		case *logrus.Fields:
			fields = *arg
		case LogInfo:
			fields = newFieldsWithLogInfo(fields, arg)
		case Error:
			lvl = getLogLevel(arg.Code)
			fields[LogKeyStatusCode] = arg.StatusCode()
			others = append(others, arg.Error())

		case *Error:
			lvl = getLogLevel(arg.Code)
			others = append(others, arg.Error())

		default:
			others = append(others, arg)
		}
	}
	return fields, lvl, others
}

func newFieldsWithLogInfo(exist logrus.Fields, val LogInfo) logrus.Fields {
	rs := exist
	for k := range val {
		rs[k] = val[k]
	}

	return rs
}

func (l *logSimple) Errorf(str string, vals ...interface{}) error {
	l.log.Errorf(str, vals...)
	return nil
}
