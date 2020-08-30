package gerr

import (
	"bytes"
	"strconv"
)

// Error error structure
//
// TracingID: id for tracing session
// Code: status http code
// Message: error message
// Target: object is mentioned in message
// Errors: error details
type Error struct {
	TracingID string
	Code      int
	Message   string
	Target    string
	Errors    []*Error
}

func (e Error) Error() string {
	b := new(bytes.Buffer)

	if e.TracingID != "" {
		pad(b, "tracingId: ")
		b.WriteString(e.TracingID)
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
	return e.Code
}

// ToResponseError make response err
func (e Error) ToResponseError() ErrResponse {
	return NewResponseError(e)
}
