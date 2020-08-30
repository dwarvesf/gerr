package gerr

// ErrDetailResponse error detail
type ErrDetailResponse map[string]interface{}

// ErrResponse error presentation
type ErrResponse struct {
	Message string            `json:"message,omitempty"`
	Errors  ErrDetailResponse `json:"errors,omitempty"`
}

// NewResponseError make err response from system Error
func NewResponseError(err Error) ErrResponse {
	return doMakeErrResponse(err)
}

func doMakeErrResponse(err Error) ErrResponse {
	return ErrResponse{
		Message: err.Message,
		// Errors:    doMakeErrDetails(err.Errors),
		Errors: normalizeResult(doMakeErrDetails(err.Errors)),
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
		tmp := doMakeErrDetail(itm)
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
