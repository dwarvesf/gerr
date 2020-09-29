package gerr

import (
	"fmt"
	"log"
	"runtime"
)

// ErrorBuilder ..
type ErrorBuilder struct {
	TraceID string
	Code    int
	Message string
	Errors  []*Error
}

// Err make a error
func (e ErrorBuilder) Err(args ...interface{}) Error {
	newArg := []interface{}{
		e.Code,
		e.Message,
		e.Errors,
		skipCaller(1),
	}

	newArg = append(
		newArg, args...,
	)
	return E(
		newArg...,
	)
}

// New make a error builder
func New(args ...interface{}) ErrorBuilder {
	e := ErrorBuilder{}
	for _, arg := range args {
		switch arg := arg.(type) {

		case Message:
			e.Message = string(arg)

		case string:
			e.Message = arg

		case int:
			e.Code = arg

		case Code:
			e.Code = int(arg)

		case *Error:
			// Make a copy
			copy := arg
			e.Errors = append(e.Errors, copy)

		case Error:
			// Make a copy
			copy := arg
			e.Errors = append(e.Errors, &copy)

		case []Error:
			// Make a copy
			for idx := range arg {
				currErr := arg[idx]

				e.Errors = append(e.Errors, &currErr)
			}

		case []*Error:
			// Make a copy
			copy := arg
			e.Errors = append(e.Errors, copy...)

		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			msg := fmt.Sprintf("unknown type %T, value %v in error call", arg, arg)
			errChild := Error{Message: msg}
			e.Errors = append(e.Errors, &errChild)
		}
	}

	return e
}
