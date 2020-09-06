package gerr

import (
	"fmt"
	"log"
	"runtime"
)

// CombinedError combined key error model
type CombinedError struct {
	Code    int
	Message string
	Target  string
	Items   []CombinedItem
}

// CombinedItem detail for combined key error model
// Keys should be string array
//     - ["items", "0", "productId"]
// Message is error message
type CombinedItem struct {
	Keys    []string
	Message string
}

// ToError make error from combined key error
func (e CombinedError) ToError() *Error {
	return makeErrorFromCombinedError(e)
}

// CombinedE helper func for init combined key error
func CombinedE(args ...interface{}) CombinedError {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}

	e := CombinedError{}
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

		case *CombinedItem:
			// Make a copy
			copy := *arg
			e.Items = append(e.Items, copy)

		case CombinedItem:
			// Make a copy
			copy := arg
			e.Items = append(e.Items, copy)

		case []CombinedItem:
			// Make a copy
			copy := arg
			e.Items = append(e.Items, copy...)

		case []*CombinedItem:
			// Make a copy
			for idx := range arg {
				currErr := arg[idx]

				e.Items = append(e.Items, *currErr)
			}

		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			msg := fmt.Sprintf("unknown type %T, value %v in error call", arg, arg)
			errChild := CombinedItem{Message: msg}
			e.Items = append(e.Items, errChild)
		}
	}
	return e
}

// makeErrorFromCombinedError make error form combined key error
func makeErrorFromCombinedError(err CombinedError) *Error {
	rs := &Error{
		Code:    err.Code,
		Message: err.Message,
		Target:  err.Target,
	}

	for idx := range err.Items {
		itm := err.Items[idx]
		doMakeChildren(itm.Keys, itm.Message, rs)
	}
	return rs
}

func doMakeChildren(keys []string, msg string, err *Error) {
	if len(keys) == 1 {
		newNode := &Error{Target: keys[0], Message: msg}
		if len(err.Errors) > 0 {
			err.Errors = append(err.Errors, newNode)
		} else {
			err.Errors = []*Error{newNode}
		}
		return
	}

	currNode := err
	currKey, keys := popKey(keys)
	if currKey == "" {
		return
	}

	found := false
	for idx := range currNode.Errors {
		childNode := currNode.Errors[idx]
		childKey := childNode.Target
		if currKey == childKey {
			found = true
			currNode = childNode
			break
		}
	}

	if !found {
		newNode := &Error{Target: currKey}
		if len(currNode.Errors) > 0 {
			currNode.Errors = append(currNode.Errors, newNode)
		} else {
			currNode.Errors = []*Error{newNode}
		}
		currNode = newNode
	}
	doMakeChildren(keys, msg, currNode)
}

func popKey(arr []string) (string, []string) {
	if len(arr) <= 0 {
		return "", arr
	}

	element := arr[0]
	if len(arr) > 1 {
		arr = arr[1:]
	} else {
		arr = make([]string, 0)
	}
	return element, arr
}
