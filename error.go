package gerr

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Error error structure
//
// TraceID: id for tracing session
// Code: status http code
// Message: error message
// Target: object is mentioned in message
// Errors: error details
type Error struct {
	TraceID string
	Code    int
	Message string
	Target  string
	Op      string
	Errors  []*Error
	trace   *stacktrace
}

func (e Error) doError() string {
	b := new(bytes.Buffer)

	if e.TraceID != "" {
		pad(b, "traceId: ")
		b.WriteString(e.TraceID)
	}

	if e.Op != "" {
		pad(b, "op: ")
		b.WriteString(e.Op)
	}

	if e.Code > 0 {
		pad(b, "code: ")
		b.WriteString(strconv.Itoa(e.Code))
	}

	if e.Target != "" {
		pad(b, "target: ")
		b.WriteString(e.Target)
	}

	if e.Message != "" {
		pad(b, "message: ")
		b.WriteString(e.Message)
	}

	if e.trace != nil {
		pad(b, "trace: \n")
		b.WriteString(e.trace.file + ":" + strconv.Itoa(e.trace.line) + " (" + e.trace.function + ")")
		if e.trace.fullTrace != "" {
			b.WriteString("\n" + e.trace.fullTrace)
		}
	}

	if e.Errors != nil {
		// Indent on new line if we are cascading non errors.
		for idx := range e.Errors {
			itm := e.Errors[idx]
			pad(b, Separator)
			b.WriteString(itm.Error())

		}
	}

	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

func (e Error) Error() string {
	return e.formatFull()
}

// StatusCode status code in Error
//
// http status code
func (e Error) StatusCode() int {
	return getStatusCode(e.Code)
}

// ToResponseError make response err
func (e Error) ToResponseError(traceID string) ErrResponse {
	e.TraceID = traceID
	return NewResponseError(e)
}

func (e Error) formatFull() string {
	var str string
	newline := func() {
		if str != "" && !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
	}

	curr := e
	str += curr.Message
	trace := curr.trace

	if trace != nil {
		newline()
		if trace.function != "" {
			str += fmt.Sprintf(" --- at %v:%v (%v) ---", trace.file, trace.line, trace.function)
		} else {
			str += fmt.Sprintf(" --- at %v:%v ---", trace.file, trace.line)
		}
	}

	for idx := range curr.Errors {
		newline()
		itm := curr.Errors[idx]
		str += itm.formatFull()
	}
	return str
}

func (e Error) formatBrief() string {
	var str string
	concat := func(msg string) {
		if str != "" && msg != "" {
			str += ": "
		}
		str += msg
	}

	curr := e
	concat(curr.Message)
	for idx := range curr.Errors {
		itm := curr.Errors[idx]
		concat(itm.formatBrief())
	}
	return str
}
