package gerr

import "net/http"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = E("internal server error", http.StatusInternalServerError)

	// ErrServiceUnavailable will throw if the requested service is not available
	ErrServiceUnavailable = E("service unavailable", http.StatusServiceUnavailable)

	// ErrNotFound will throw if the requested item is not found
	ErrNotFound = E("not found", http.StatusNotFound)

	// ErrBadRequest will throw if the request is invalid
	ErrBadRequest = E("bad request", http.StatusBadRequest)

	// ErrUnauthorized will throw if client failed to authenticate with the server
	ErrUnauthorized = E("unauthorized", http.StatusUnauthorized)

	// ErrForbidden will throw if client authenticated but does not have permission to access the requested resource
	ErrForbidden = E("forbidden", http.StatusForbidden)
)
