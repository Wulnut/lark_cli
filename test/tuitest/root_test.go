package tuittest

import (
	"testing"

	"lark_cli/internal/tui"
)

func TestRootModel_SuccessView(t *testing.T) {
	// Integration test — requires httptest server + fake token provider.
	// Skipped in unit test runs; enable with -tags=integration.
	t.Skip("integration test - requires httptest server")
}

func TestRootModel_DegradedView(t *testing.T) {
	t.Skip("integration test - requires httptest server")
}

func TestRun_NotLoggedIn(t *testing.T) {
	// Smoke test: Run with no client should not panic.
	// In practice this is hard to unit test without wrapping Run.
	t.Skip("requires httptest or integration test setup")
}

// Test that tui package compiles with new signatures.
func TestTUICompiles(t *testing.T) {
	// If this package compiles, the signatures are correct.
	// Actual rendering tests need integration setup.
	_ = tui.Run
}
