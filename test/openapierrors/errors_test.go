package openapierrors

import (
	"net/http"
	"strings"
	"testing"

	openapierrors "lark_cli/internal/openapi/errors"
)

func TestLookupPolicy_KnownCodes(t *testing.T) {
	tests := []struct {
		code             int
		wantCategory     openapierrors.Category
		wantRetryable    bool
		wantRefreshToken bool
		wantMaxRetry     int
		wantHTTPStatus   int
	}{
		// Auth
		{10021, openapierrors.CategoryAuth, true, true, 1, 401},
		{10022, openapierrors.CategoryAuth, true, true, 1, 401},
		{10301, openapierrors.CategoryAuth, true, true, 1, 401},
		{10302, openapierrors.CategoryAuth, false, false, 0, 401},
		// Permission
		{10001, openapierrors.CategoryPermission, false, false, 0, 403},
		{10002, openapierrors.CategoryPermission, false, false, 0, 403},
		{10004, openapierrors.CategoryPermission, false, false, 0, 403},
		{10210, openapierrors.CategoryPermission, false, false, 0, 403},
		{10211, openapierrors.CategoryPermission, false, false, 0, 403},
		{10404, openapierrors.CategoryPermission, false, false, 0, 403},
		// Rate limit
		{10429, openapierrors.CategoryRateLimit, true, false, 3, 429},
		{10430, openapierrors.CategoryRateLimit, false, false, 0, 429},
		// Validation
		{20001, openapierrors.CategoryValidation, false, false, 0, 400},
		{20005, openapierrors.CategoryValidation, false, false, 0, 400},
		{20006, openapierrors.CategoryValidation, false, false, 0, 400},
		{20090, openapierrors.CategoryValidation, false, false, 0, 400},
		// Not found
		{13001, openapierrors.CategoryNotFound, false, false, 0, 404},
		{30005, openapierrors.CategoryNotFound, false, false, 0, 404},
		{30006, openapierrors.CategoryNotFound, false, false, 0, 404},
		// Server
		{50006, openapierrors.CategoryServer, true, false, 2, 500},
		// Other
		{9999, openapierrors.CategoryClient, false, false, 0, 400},
		{1000051942, openapierrors.CategoryClient, false, false, 0, 400},
	}

	for _, tt := range tests {
		name := strings.TrimPrefix(t.Name(), "TestLookupPolicy_")
		t.Run(http.StatusText(tt.wantHTTPStatus)+"/"+string(tt.wantCategory)+"/"+name, func(t *testing.T) {
			p := openapierrors.LookupPolicy(tt.code)
			if p.Category != tt.wantCategory {
				t.Errorf("code %d: Category = %q, want %q", tt.code, p.Category, tt.wantCategory)
			}
			if p.Retryable != tt.wantRetryable {
				t.Errorf("code %d: Retryable = %v, want %v", tt.code, p.Retryable, tt.wantRetryable)
			}
			if p.RefreshToken != tt.wantRefreshToken {
				t.Errorf("code %d: RefreshToken = %v, want %v", tt.code, p.RefreshToken, tt.wantRefreshToken)
			}
			if p.MaxRetry != tt.wantMaxRetry {
				t.Errorf("code %d: MaxRetry = %d, want %d", tt.code, p.MaxRetry, tt.wantMaxRetry)
			}
			if p.HTTPStatus != tt.wantHTTPStatus {
				t.Errorf("code %d: HTTPStatus = %d, want %d", tt.code, p.HTTPStatus, tt.wantHTTPStatus)
			}
		})
	}
}

func TestLookupPolicy_UnknownCode(t *testing.T) {
	p := openapierrors.LookupPolicy(99999)
	if p.Category != openapierrors.CategoryUnknown {
		t.Errorf("unknown code: Category = %q, want unknown", p.Category)
	}
	if p.Retryable {
		t.Errorf("unknown code: Retryable should be false")
	}
	if p.HTTPStatus != 0 {
		t.Errorf("unknown code: HTTPStatus = %d, want 0", p.HTTPStatus)
	}
	if p.UserHint == "" {
		t.Error("unknown code: should have a UserHint")
	}
}

func TestNormalizeError_Success(t *testing.T) {
	err := openapierrors.NormalizeError(http.StatusOK, []byte(`{"data":{}}`), 0, "")
	if err != nil {
		t.Errorf("NormalizeError on success: expected nil, got %v", err)
	}
}

func TestNormalizeError_WithCode(t *testing.T) {
	body := []byte(`{"err_code":10022,"err_msg":"Check Token Failed"}`)
	err := openapierrors.NormalizeError(http.StatusUnauthorized, body, 10022, "Check Token Failed")

	openErr, ok := err.(*openapierrors.OpenAPIError)
	if !ok {
		t.Fatalf("expected *OpenAPIError, got %T", err)
	}
	if openErr.Code != 10022 {
		t.Errorf("Code = %d, want 10022", openErr.Code)
	}
	if openErr.HTTPStatus != http.StatusUnauthorized {
		t.Errorf("HTTPStatus = %d, want %d", openErr.HTTPStatus, http.StatusUnauthorized)
	}
	if openErr.Category != openapierrors.CategoryAuth {
		t.Errorf("Category = %q, want auth", openErr.Category)
	}
	if openErr.UserHint == "" {
		t.Error("UserHint should not be empty")
	}
}

func TestNormalizeError_Non2xxWithoutCode(t *testing.T) {
	body := []byte("service unavailable")
	err := openapierrors.NormalizeError(http.StatusServiceUnavailable, body, 0, "http status 503")

	openErr, ok := err.(*openapierrors.OpenAPIError)
	if !ok {
		t.Fatalf("expected *OpenAPIError, got %T", err)
	}
	if openErr.HTTPStatus != http.StatusServiceUnavailable {
		t.Errorf("HTTPStatus = %d, want %d", openErr.HTTPStatus, http.StatusServiceUnavailable)
	}
	if openErr.Category != openapierrors.CategoryUnknown {
		t.Errorf("Category = %q, want unknown", openErr.Category)
	}
}

func TestOpenAPIError_ErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      *openapierrors.OpenAPIError
		contains []string
	}{
		{
			name: "with code",
			err: &openapierrors.OpenAPIError{
				HTTPStatus: 401,
				Code:       10022,
				Message:    "Check Token Failed",
			},
			contains: []string{"10022", "Check Token Failed"},
		},
		{
			name: "without code with message",
			err: &openapierrors.OpenAPIError{
				HTTPStatus: 503,
				Code:       0,
				Message:    "http status 503: upstream unavailable",
			},
			contains: []string{"503", "upstream unavailable"},
		},
		{
			name: "without code without message",
			err: &openapierrors.OpenAPIError{
				HTTPStatus: 500,
				Code:       0,
				Message:    "",
			},
			contains: []string{"500"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.err.Error()
			for _, sub := range tt.contains {
				if !strings.Contains(s, sub) {
					t.Errorf("Error() = %q, want to contain %q", s, sub)
				}
			}
		})
	}
}

func TestPoliciesMapNotEmpty(t *testing.T) {
	// Smoke test: LookupPolicy should return valid policies for known codes.
	p := openapierrors.LookupPolicy(10021)
	if p.Code != 10021 {
		t.Error("Policies map appears empty or LookupPolicy is broken")
	}
}
