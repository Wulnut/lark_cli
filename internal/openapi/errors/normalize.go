package errors

// NormalizeError converts raw API response data into a normalized OpenAPIError.
func NormalizeError(httpStatus int, body []byte, code int, msg string) error {
	// Success case
	if code == 0 && httpStatus >= 200 && httpStatus < 300 {
		return nil
	}

	policy := LookupPolicy(code)

	return &OpenAPIError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    msg,
		Category:   policy.Category,
		UserHint:   policy.UserHint,
		DevHint:    policy.DevHint,
		RawBody:    body,
	}
}

// NormalizeFromPolicy creates an OpenAPIError from an already-resolved Policy.
// Use this when the caller already has the policy (e.g., from LookupPolicy).
func NormalizeFromPolicy(httpStatus int, body []byte, code int, msg string, policy Policy) *OpenAPIError {
	return &OpenAPIError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    msg,
		Category:   policy.Category,
		UserHint:   policy.UserHint,
		DevHint:    policy.DevHint,
		RawBody:    body,
	}
}
