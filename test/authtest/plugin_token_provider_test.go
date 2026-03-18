package authtest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"lark_cli/internal/auth"
	"lark_cli/internal/session"
)

// fakeClock is a controllable clock for testing.
type fakeClock struct {
	now time.Time
}

func (c *fakeClock) Now() time.Time          { return c.now }
func (c *fakeClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

func newTestServer(t *testing.T, token string, expireTime int, errorCode int, errorMsg string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"data": map[string]any{
				"token":       token,
				"expire_time": expireTime,
			},
			"error": map[string]any{
				"code": errorCode,
				"msg":  errorMsg,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestPluginTokenProvider_FirstGetFetchesAndPersists(t *testing.T) {
	server := newTestServer(t, "p-test-token-123456", 7200, 0, "success")
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: now}

	// Save a session with user_key but no token
	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		UserKey:   "ou_test",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := store.Save(ctx, sess); err != nil {
		t.Fatalf("Save: %v", err)
	}

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(store, client, 10*time.Minute, clock)

	token, err := provider.Get(ctx)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if token != "p-test-token-123456" {
		t.Errorf("token = %q", token)
	}

	// Verify session was updated
	loaded, err := store.Load(ctx)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.PluginAccessToken != "p-test-token-123456" {
		t.Errorf("persisted token = %q", loaded.PluginAccessToken)
	}
}

func TestPluginTokenProvider_CachedTokenReused(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := map[string]any{
			"data":  map[string]any{"token": "p-cached", "expire_time": 7200},
			"error": map[string]any{"code": 0, "msg": "success"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: now}

	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		UserKey:   "ou_test",
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = store.Save(ctx, sess)

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(store, client, 10*time.Minute, clock)

	// First call fetches
	_, _ = provider.Get(ctx)
	if callCount != 1 {
		t.Fatalf("expected 1 fetch, got %d", callCount)
	}

	// Second call within expiry should use cache (no additional fetch)
	clock.Advance(30 * time.Minute) // Still within 2h - 10min leeway
	_, _ = provider.Get(ctx)
	if callCount != 1 {
		t.Errorf("expected still 1 fetch (cached), got %d", callCount)
	}
}

func TestPluginTokenProvider_ExpiredTokenRefetches(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := map[string]any{
			"data":  map[string]any{"token": "p-refreshed", "expire_time": 7200},
			"error": map[string]any{"code": 0, "msg": "success"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: now}

	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		UserKey:   "ou_test",
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = store.Save(ctx, sess)

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(store, client, 10*time.Minute, clock)

	// First call
	_, _ = provider.Get(ctx)
	if callCount != 1 {
		t.Fatalf("expected 1 fetch, got %d", callCount)
	}

	// Advance past expiry (2h total, leeway 10min → expired at 1h50m)
	clock.Advance(2 * time.Hour)
	_, _ = provider.Get(ctx)
	if callCount != 2 {
		t.Errorf("expected 2 fetches after expiry, got %d", callCount)
	}
}

func TestPluginTokenProvider_ForceRefreshAlwaysFetches(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := map[string]any{
			"data":  map[string]any{"token": "p-forced", "expire_time": 7200},
			"error": map[string]any{"code": 0, "msg": "success"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: now}

	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		UserKey:   "ou_test",
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = store.Save(ctx, sess)

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(store, client, 10*time.Minute, clock)

	// Get once (cached)
	_, _ = provider.Get(ctx)

	// ForceRefresh should always fetch
	token, err := provider.ForceRefresh(ctx)
	if err != nil {
		t.Fatalf("ForceRefresh: %v", err)
	}
	if token != "p-forced" {
		t.Errorf("token = %q", token)
	}
	if callCount != 2 {
		t.Errorf("expected 2 fetches, got %d", callCount)
	}
}

func TestPluginTokenProvider_NotLoggedIn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	client := auth.NewPluginTokenClient(http.DefaultClient, "http://localhost", "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(store, client, 10*time.Minute, &fakeClock{now: time.Now()})

	_, err := provider.Get(ctx)
	if err == nil {
		t.Fatal("expected error for missing session")
	}
}

func TestHeaderProvider_ReturnsHeaders(t *testing.T) {
	server := newTestServer(t, "p-header-token", 7200, 0, "success")
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: now}

	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		UserKey:   "ou_header_test",
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = store.Save(ctx, sess)

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	tokenProvider := auth.NewPluginTokenProviderWithClock(store, client, 10*time.Minute, clock)
	headerProvider := auth.NewHeaderProvider(store, tokenProvider)

	headers, err := headerProvider.Headers(ctx)
	if err != nil {
		t.Fatalf("Headers: %v", err)
	}
	if headers["X-Plugin-Token"] != "p-header-token" {
		t.Errorf("X-Plugin-Token = %q", headers["X-Plugin-Token"])
	}
	if headers["X-User-Key"] != "ou_header_test" {
		t.Errorf("X-User-Key = %q", headers["X-User-Key"])
	}
}

func TestHeaderProvider_NotLoggedIn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	client := auth.NewPluginTokenClient(http.DefaultClient, "http://localhost", "pid", "psecret")
	tokenProvider := auth.NewPluginTokenProviderWithClock(store, client, 10*time.Minute, &fakeClock{now: time.Now()})
	headerProvider := auth.NewHeaderProvider(store, tokenProvider)

	_, err := headerProvider.Headers(ctx)
	if err == nil {
		t.Fatal("expected error for missing session")
	}
}

func TestMask(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "***"},
		{"abc", "***"},
		{"abcdef", "***"},
		{"abcdefg", "abcdef***"},
		{"p-49257489-f7d7-4cd6-b34f", "p-4925***"},
	}
	for _, tt := range tests {
		got := auth.Mask(tt.input)
		if got != tt.want {
			t.Errorf("Mask(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
