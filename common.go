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

func getStatusCode(code int) int {
	if code <= httpMaxLength {
		return code
	}

	if code < internalCodeMax {
		return http.StatusInternalServerError
	}

	if code < serviceCodeMax {
		return http.StatusInternalServerError
	}

	if code < businessCodeMax {
		return http.StatusBadRequest
	}

	return http.StatusBadRequest
}
