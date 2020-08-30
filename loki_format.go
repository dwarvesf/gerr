package gerr

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/sirupsen/logrus"
)

const (
	defaultTimestampFormat = time.RFC3339
)

type fieldKey string

// FieldMap allows customization of the key names for default logrus.fields.
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

// lokiFormatter formats logs into parsable json
type lokiFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// Force quoting of all values
	ForceQuote bool

	// DisableQuote disables quoting for all values.
	// DisableQuote will have a lower priority than ForceQuote.
	// If both of them are set to true, quote will be forced on all values.
	DisableQuote bool

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// FieldMap allows users to customize the names of keys for default logrus.fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 logrus.FieldKeyTime:  "@timestamp",
	// 		 logrus.FieldKeyLevel: "@level",
	// 		 logrus.FieldKeyMsg:   "@message",
	// 		 logrus.FieldKeyFunc:  "@caller",
	//    },
	// }
	FieldMap FieldMap

	// PrettyPrint will indent all json logs
	PrettyPrint bool

	terminalInitOnce sync.Once

	// The max length of the level text, generated dynamically on init
	levelTextMaxLength int
}

// Format renders a single log entry
func (f *lokiFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	var funcVal, fileVal string

	fixedKeys := make([]string, 0, 4+len(data))
	if !f.DisableTimestamp {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyTime))
	}
	fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyLevel))

	if entry.HasCaller() {
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		} else {
			funcVal = entry.Caller.Function
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}

		if funcVal != "" {
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyFunc))
		}
		if fileVal != "" {
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyFile))
		}
	}

	if !f.DisableSorting {
		sort.Strings(keys)
		fixedKeys = append(fixedKeys, keys...)
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	f.terminalInitOnce.Do(func() { f.init(entry) })

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	b.WriteByte('{')

	numOfKey := len(fixedKeys)
	for idx := range fixedKeys {
		key := fixedKeys[idx]
		var value interface{}
		switch {
		case key == f.FieldMap.resolve(logrus.FieldKeyTime):
			value = entry.Time.Format(timestampFormat)
		case key == f.FieldMap.resolve(logrus.FieldKeyLevel):
			value = entry.Level.String()
		case key == f.FieldMap.resolve(logrus.FieldKeyMsg):
			value = entry.Message
		case key == f.FieldMap.resolve(logrus.FieldKeyFunc) && entry.HasCaller():
			value = funcVal
		case key == f.FieldMap.resolve(logrus.FieldKeyFile) && entry.HasCaller():
			value = fileVal
		default:
			value = data[key]
		}
		f.appendKeyValue(b, key, value)
		if idx < numOfKey-1 {
			b.WriteByte(',')
		}
	}

	b.WriteString(" } ")
	b.WriteString(entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *lokiFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *lokiFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func (f *lokiFormatter) needsQuoting(text string) bool {
	if f.ForceQuote {
		return true
	}

	if f.DisableQuote {
		return false
	}

	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *lokiFormatter) init(entry *logrus.Entry) {
	// Get the max length of the level text
	for _, level := range logrus.AllLevels {
		levelTextLength := utf8.RuneCount([]byte(level.String()))
		if levelTextLength > f.levelTextMaxLength {
			f.levelTextMaxLength = levelTextLength
		}
	}
}
