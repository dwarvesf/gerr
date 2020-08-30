package gerr

import (
	"fmt"
	"log"
	"runtime"
)

// Target target for an error
//
// string
type Target string

// Message message for an error
//
// string
type Message string

// Code status code for an error
//
// int
type Code int

// TracingID for an error
//
// string
type TracingID string

// E builds an error value from its arguments.
func E(args ...interface{}) Error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}

	e := Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Target:
			e.Target = string(arg)

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
			copy := &arg
			e.Errors = append(e.Errors, copy)

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
