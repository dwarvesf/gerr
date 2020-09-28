package gerr

const (
	serviceCodeMin = iota + serviceCodeMax

	// ErrSvcTimeout sevice Timeout
	ErrSvcTimeout

	// ErrSvcLostConnection sevice Lost Connection
	ErrSvcLostConnection

	// ErrSvcReconnectTimeOut sevice Reconnect TimeOut
	ErrSvcReconnectTimeOut

	// ErrSvcAuthRequired sevice Auth Required
	ErrSvcAuthRequired

	// ErrSvcPermissionRequired sevice Permission Required
	ErrSvcPermissionRequired

	// NOTE: all service code should be add above serviceCodeLength
	serviceCodeLength
)

const (
	// ServiceCodeCustomStart custom Service code index start
	ServiceCodeCustomStart = serviceCodeLength
	serviceCodeMax         = int(20000)
)

var serviceMsg = map[int]string{
	ErrSvcTimeout:            "timeout",
	ErrSvcLostConnection:     "lost connection",
	ErrSvcReconnectTimeOut:   "reconnect timeOut",
	ErrSvcAuthRequired:       "auth required",
	ErrSvcPermissionRequired: "permission required",
}

func getServiceMessage(code int) string {
	if msg, ok := serviceMsg[code]; ok {
		return msg
	}
	return ""
}
