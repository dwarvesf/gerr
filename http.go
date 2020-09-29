package gerr

// ErrDetailResponse error detail
type ErrDetailResponse map[string]interface{}

// ErrResponse error presentation
type ErrResponse struct {
	Message string            `json:"message,omitempty"`
	TraceID string            `json:"trace_id,omitempty"`
	Code    int               `json:"code,omitempty"`
	Errors  ErrDetailResponse `json:"errors,omitempty"`
}

// NewResponseError make err response from system Error
func NewResponseError(err Error) ErrResponse {
	return doMakeErrResponse(err)
}

func doMakeErrResponse(err Error) ErrResponse {
	details := []*Error{}

	for idx := range err.Errors {
		itm := err.Errors[idx]
		if dt := castError(itm); dt != nil {
			details = append(details, dt)
		}
	}
	return ErrResponse{
		Message: err.Message,
		TraceID: err.TraceID,
		Code:    err.Code,
		// Errors:    doMakeErrDetails(err.Errors),
		Errors: normalizeResult(doMakeErrDetails(details)),
	}
}

func normalizeResult(m map[string]interface{}) map[string]interface{} {
	rs := map[string]interface{}{}

	for k := range m {
		val := nomalizeEmptyKey(m[k])
		currMap, ok := val.(map[string]interface{})
		if !ok {
			rs[k] = val
			continue
		}
		rs[k] = normalizeResult(currMap)
	}

	return rs
}

func nomalizeEmptyKey(val interface{}) interface{} {
	str, ok := val.(string)
	if ok {
		return []interface{}{str}
	}

	m, ok := val.(map[string]interface{})
	if !ok {
		return val
	}

	if len(m) == 1 {
		dt, ok := m[""]
		if ok {
			return dt
		}
	}
	return val
}

func doMakeErrDetails(errs []*Error) map[string]interface{} {
	rs := map[string]interface{}{}

	for idx := range errs {
		itm := errs[idx]
		tmp := doMakeErrDetail(itm)
		for key := range tmp {
			rs[key] = tmp[key]
		}
	}

	return rs
}

func castError(err error) *Error {
	if dt, ok := err.(*Error); ok {
		return dt
	}

	if dt, ok := err.(Error); ok {
		return &dt
	}

	return nil
}

func doMakeErrDetail(err *Error) map[string]interface{} {
	if !hasChildren(*err) {
		rs := map[string]interface{}{}
		k := err.Target
		v := err.Message
		rs[k] = v
		return rs
	}

	rs := map[string]interface{}{}

	for idx := range err.Errors {
		itm := err.Errors[idx]
		errItm := castError(itm)
		if errItm == nil {
			continue
		}
		tmp := doMakeErrDetail(errItm)
		for key := range tmp {
			tmpData, ok := rs[key]
			if !ok {
				rs[key] = tmp[key]
				continue
			}

			parsedItmData, parsedItmOk := tmpData.(string)
			if parsedItmOk {
				rs[key] = []interface{}{parsedItmData, tmp[key]}
			}

			parsedListData, parsedListOk := tmpData.([]interface{})
			if parsedListOk {
				rs[key] = append(parsedListData, tmp[key])
				continue
			}
		}
	}
	return map[string]interface{}{err.Target: rs}
}

func hasChildren(err Error) bool {
	return len(err.Errors) > 0
}
