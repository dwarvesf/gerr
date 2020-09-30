package gerr

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/dwarvesf/gerr/cleanpath"
)

/*
CleanPath function is applied to file paths before adding them to a stacktrace.
By default, it makes the path relative to the $GOPATH environment variable.
To remove some additional prefix like "github.com" from file paths in
stacktraces, use something like:
	stacktrace.CleanPath = func(path string) string {
		path = cleanpath.RemoveGoPath(path)
		path = strings.TrimPrefix(path, "github.com/")
		return path
	}
*/
var CleanPath = cleanpath.RemoveGoPath

// SetCleanPath set clean path func
func SetCleanPath(path string) {
	CleanPath = cleanpath.RemoveProjectPath(path)
}

// SetCleanPathFunc set clean path func
func SetCleanPathFunc(fn ...func(string) string) {
	chain := append(
		[]func(string) string{
			cleanpath.RemoveGoPath,
		},
		fn...,
	)
	CleanPath = cleanpath.PathHandleChain(chain...)
}

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

// TraceID for an error
//
// string
type TraceID string

// Op function name
//
// string
type Op string

// TraceError data for an error
//
// string
type TraceError struct {
	Error error
}

type skipCaller int

// E builds an error value from its arguments.
func E(args ...interface{}) Error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	skip := 1

	e := Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Target:
			e.Target = string(arg)

		case Message:
			e.Message = string(arg)

		case Op:
			e.Op = string(arg)

		case string:
			e.Message = arg

		case int:
			e.Code = arg

		case Code:
			e.Code = int(arg)

		case *Error:
			// Make a copy
			copy := arg
			if e.Code <= 0 {
				e = *copy
				break
			}
			e.Errors = append(e.Errors, copy)

		case Error:
			// Make a copy
			copy := arg
			if e.Code <= 0 {
				e = copy
				break
			}
			e.Errors = append(e.Errors, &copy)

		case error:
			// Make a copy
			if itm := castError(arg); itm != nil {
				e.Errors = append(e.Errors, itm)
			} else {
				e.Errors = append(e.Errors, arg)
			}

		case []error:
			// Make a copy
			copy := arg
			e.Errors = append(e.Errors, copy...)

		case []Error:
			// Make a copy
			errs := []error{}
			for idx := range arg {
				errs = append(errs, &arg[idx])
			}

			e.Errors = append(e.Errors, errs...)

		case []*Error:
			errs := []error{}
			for idx := range arg {
				errs = append(errs, arg[idx])
			}

			e.Errors = append(e.Errors, errs...)

		case TraceError:
			e.Trace = arg.Error

		case skipCaller:
			skip = skip + int(arg)

		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			msg := fmt.Sprintf("unknown type %T, value %v in error call", arg, arg)
			errChild := Error{Message: msg}
			e.Errors = append(e.Errors, &errChild)
		}
	}

	if e.Message == "" && e.Code > 0 {
		e.Message = getDefaultMessage(e.Code)
	}

	if e.Op == "" {
		op := getStackTrace(skip)
		e.Op = op.function
		e.trace = &op
	}

	return e
}

func getStackTrace(skip int) stacktrace {
	skip = skip + 1
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)
	stackStr := string(buf[0:stackSize])

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return newStackTrace("", 0, "?", stackStr)
	}

	if CleanPath != nil {
		file = CleanPath(file)
	}

	fn := runtime.FuncForPC(pc)
	fnName := ""
	if fn != nil {

		fnTemp := shortFuncName(fn)
		if fnTemp != "" {
			fnName = fnTemp
		}

		longName := fn.Name()
		stackStr = stackStr[strings.LastIndex(stackStr, longName):]
		stackStrList := strings.Split(stackStr, "\n\n")
		stackStr = stackStrList[0]
		if fnTemp == "init" {
			stackStr = ""
			fnName = ""
		}
	}

	return newStackTrace(file, line, fnName, stackStr)
}

/* "FuncName" or "Receiver.MethodName" */
func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/palantir/shield/package.FuncName"
	// - "github.com/palantir/shield/package.Receiver.MethodName"
	// - "github.com/palantir/shield/package.(*PtrReceiver).MethodName"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}

// Et builds an error value from its arguments.
// NOTE: MUST add skipCaller(1) for wrapper function
func Et(target string, args ...interface{}) Error {
	t := Target(target)
	args = append(args, t, skipCaller(1))
	return E(args...)
}

// Trace make trace error data for error
func Trace(err error) TraceError {
	return TraceError{err}
}
