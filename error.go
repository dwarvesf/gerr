package gerr

import (
	"bytes"
	"strconv"
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

func (e Error) Error() string {
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

// StatusCode status code in Error
//
// http status code
func (e Error) StatusCode() int {
	return getStatusCode(e.Code)
}

// ToResponseError make response err
func (e Error) ToResponseError() ErrResponse {
	return NewResponseError(e)
}
