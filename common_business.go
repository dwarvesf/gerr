package gerr

const (
	businessCodeMin = iota + serviceCodeMax

	// ErrAuthWrongCredential username or password is incorrect
	ErrAuthWrongCredential

	// ErrAuthNoPermission no permission
	ErrAuthNoPermission

	// ErrAuthTokenInvalid token invalid
	ErrAuthTokenInvalid

	// ErrAuthTokenExpired token expired
	ErrAuthTokenExpired

	// ErrRecordNotFound record not found
	ErrRecordNotFound

	// ErrIDInvalid id invalid
	ErrIDInvalid

	// NOTE: all business code should be add above businessCodeLength
	businessCodeLength
)

const (
	// BusinessCodeCustomStart custom business code index start
	BusinessCodeCustomStart = businessCodeLength
)

const (
	businessCodeMax = int(999999999)
)

var businessMsg = map[int]string{
	ErrAuthWrongCredential: "username or password is incorrect",
	ErrAuthNoPermission:    "no permission",
	ErrAuthTokenInvalid:    "token invalid",
	ErrAuthTokenExpired:    "token expired",
	ErrRecordNotFound:      "record not found",
	ErrIDInvalid:           "id invalid",
}

func getBusinessMessage(code int) string {
	if msg, ok := businessMsg[code]; ok {
		return msg
	}
	return ""
}
