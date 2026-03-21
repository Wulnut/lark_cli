package errors

import "fmt"

// Category represents the error category.
type Category string

const (
	CategoryAuth       Category = "auth"
	CategoryPermission Category = "permission"
	CategoryRateLimit  Category = "rate_limit"
	CategoryServer     Category = "server"
	CategoryValidation Category = "validation"
	CategoryNotFound   Category = "not_found"
	CategoryClient     Category = "client"
	CategoryUnknown    Category = "unknown"
)

// Policy defines the handling strategy for an error code.
type Policy struct {
	Code        int
	HTTPStatus  int
	Message     string
	Category    Category
	Retryable   bool
	RefreshToken bool
	MaxRetry    int
	UserHint    string
	DevHint     string
}

// OpenAPIError represents a normalized OpenAPI error.
type OpenAPIError struct {
	HTTPStatus int
	Code       int
	Message    string
	Category   Category
	UserHint   string
	DevHint    string
	RawBody    []byte
}

func (e *OpenAPIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("openapi error: code=%d msg=%s", e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("openapi error: http_status=%d msg=%s", e.HTTPStatus, e.Message)
	}
	return fmt.Sprintf("openapi error: http_status=%d", e.HTTPStatus)
}
