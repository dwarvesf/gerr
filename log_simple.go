package gerr

import (
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
	fields, vals := detachFields(vals...)
	l.log.WithFields(fields).Info(vals...)
	return nil
}
func (l *logSimple) Debug(vals ...interface{}) error {
	fields, vals := detachFields(vals...)
	l.log.WithFields(fields).Debug(vals...)
	return nil
}
func (l *logSimple) Info(vals ...interface{}) error {
	fields, vals := detachFields(vals...)
	l.log.WithFields(fields).Info(vals...)
	return nil
}
func (l *logSimple) Warn(vals ...interface{}) error {
	fields, vals := detachFields(vals...)
	l.log.WithFields(fields).Warning(vals...)
	return nil
}
func (l *logSimple) Error(vals ...interface{}) error {
	fields, vals := detachFields(vals...)
	l.log.WithFields(fields).Warning(vals...)
	return nil
}

func detachFields(vals ...interface{}) (logrus.Fields, []interface{}) {
	var fields logrus.Fields
	others := []interface{}{}
	for idx := range vals {
		arg := vals[idx]
		switch arg := arg.(type) {
		case logrus.Fields:
			fields = arg
		case *logrus.Fields:
			fields = *arg
		case LogInfo:
			fields = newFieldsWithLogInfo(arg)
		default:
			others = append(others, arg)
		}
	}
	return fields, others
}

func newFieldsWithLogInfo(val LogInfo) logrus.Fields {
	rs := logrus.Fields{}
	for k := range val {
		rs[k] = val[k]
	}

	return rs
}

func (l *logSimple) Errorf(str string, vals ...interface{}) error {
	l.log.Errorf(str, vals...)
	return nil
}
