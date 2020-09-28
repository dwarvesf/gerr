package gerr

import "net/http"

func getDefaultMessage(code int) string {
	if code <= httpMaxLength {
		return http.StatusText(code)
	}

	if code < internalCodeMax {
		return getInternalMessage(code)
	}

	if code < serviceCodeMax {
		return getServiceMessage(code)
	}

	if code < businessCodeMax {
		return getBusinessMessage(code)
	}

	return ""
}
