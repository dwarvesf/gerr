package gerr

const (
	internalCodeMin = iota + 1000

	// ErrIOInvalidPath ..
	ErrIOInvalidPath

	// ErrIONotExist ..
	ErrIONotExist

	// ErrIOExist ..
	ErrIOExist

	// ErrIOReadFailed ..
	ErrIOReadFailed

	// ErrIOContentReachLimit ..
	ErrIOContentReachLimit

	// ErrIOWriteFailed ..
	ErrIOWriteFailed

	// NOTE: all intenal code should be add above internalCodeLength
	internalCodeLength
)

const (
	internalCodeMax = int(50000)
)

var internalMsg = map[int]string{
	ErrIOInvalidPath:       "invalid path",
	ErrIONotExist:          "not exist",
	ErrIOExist:             "exist",
	ErrIOReadFailed:        "read failed",
	ErrIOContentReachLimit: "content reach limit",
	ErrIOWriteFailed:       "write failed",
}

func getInternalMessage(code int) string {
	if msg, ok := internalMsg[code]; ok {
		return msg
	}
	return ""
}
