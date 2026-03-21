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
	"lark_cli/internal/config"
	"lark_cli/internal/session"
)

// fakeClock is a controllable clock for testing.
type fakeClock struct {
	now time.Time
}

func (c *fakeClock) Now() time.Time          { return c.now }
func (c *fakeClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

// mockConfigStore is a test double for auth.ConfigStore.
type mockConfigStore struct {
	cfg *config.Config
	err error
}

func (m *mockConfigStore) Load(ctx context.Context) (*config.Config, error) {
	return m.cfg, m.err
}

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

	cfgStore := &mockConfigStore{
		cfg: &config.Config{
			UserKey:      "ou_test",
			PluginID:     "pid",
			PluginSecret: "psecret",
			BaseURL:      server.URL,
		},
	}

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(cfgStore, store, client, 10*time.Minute, clock)

	authCtx, err := provider.GetAuthContext(ctx)
	if err != nil {
		t.Fatalf("GetAuthContext: %v", err)
	}
	if authCtx.PluginToken != "p-test-token-123456" {
		t.Errorf("token = %q", authCtx.PluginToken)
	}
	if authCtx.UserKey != "ou_test" {
		t.Errorf("userKey = %q", authCtx.UserKey)
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

	cfgStore := &mockConfigStore{
		cfg: &config.Config{
			UserKey:      "ou_test",
			PluginID:     "pid",
			PluginSecret: "psecret",
			BaseURL:      server.URL,
		},
	}

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(cfgStore, store, client, 10*time.Minute, clock)

	// First call fetches
	_, _ = provider.GetAuthContext(ctx)
	if callCount != 1 {
		t.Fatalf("expected 1 fetch, got %d", callCount)
	}

	// Second call within expiry should use cache (no additional fetch)
	clock.Advance(30 * time.Minute) // Still within 2h - 10min leeway
	_, _ = provider.GetAuthContext(ctx)
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

	cfgStore := &mockConfigStore{
		cfg: &config.Config{
			UserKey:      "ou_test",
			PluginID:     "pid",
			PluginSecret: "psecret",
			BaseURL:      server.URL,
		},
	}

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(cfgStore, store, client, 10*time.Minute, clock)

	// First call
	_, _ = provider.GetAuthContext(ctx)
	if callCount != 1 {
		t.Fatalf("expected 1 fetch, got %d", callCount)
	}

	// Advance past expiry (2h total, leeway 10min → expired at 1h50m)
	clock.Advance(2 * time.Hour)
	_, _ = provider.GetAuthContext(ctx)
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

	cfgStore := &mockConfigStore{
		cfg: &config.Config{
			UserKey:      "ou_test",
			PluginID:     "pid",
			PluginSecret: "psecret",
			BaseURL:      server.URL,
		},
	}

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(cfgStore, store, client, 10*time.Minute, clock)

	// Get once (cached)
	_, _ = provider.GetAuthContext(ctx)

	// ForceRefresh should always fetch
	authCtx, err := provider.ForceRefresh(ctx)
	if err != nil {
		t.Fatalf("ForceRefresh: %v", err)
	}
	if authCtx.PluginToken != "p-forced" {
		t.Errorf("token = %q", authCtx.PluginToken)
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

	cfgStore := &mockConfigStore{
		cfg: &config.Config{
			UserKey:      "",
			PluginID:     "pid",
			PluginSecret: "psecret",
		},
	}

	client := auth.NewPluginTokenClient(http.DefaultClient, "http://localhost", "pid", "psecret")
	provider := auth.NewPluginTokenProviderWithClock(cfgStore, store, client, 10*time.Minute, &fakeClock{now: time.Now()})

	_, err := provider.GetAuthContext(ctx)
	if err == nil {
		t.Fatal("expected error for missing user_key")
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

	cfgStore := &mockConfigStore{
		cfg: &config.Config{
			UserKey:      "ou_header_test",
			PluginID:     "pid",
			PluginSecret: "psecret",
			BaseURL:      server.URL,
		},
	}

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	tokenProvider := auth.NewPluginTokenProviderWithClock(cfgStore, store, client, 10*time.Minute, clock)
	headerProvider := auth.NewHeaderProvider(tokenProvider)

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

	cfgStore := &mockConfigStore{
		cfg: &config.Config{
			UserKey:      "",
			PluginID:     "pid",
			PluginSecret: "psecret",
		},
	}

	client := auth.NewPluginTokenClient(http.DefaultClient, "http://localhost", "pid", "psecret")
	tokenProvider := auth.NewPluginTokenProviderWithClock(cfgStore, store, client, 10*time.Minute, &fakeClock{now: time.Now()})
	headerProvider := auth.NewHeaderProvider(tokenProvider)

	_, err := headerProvider.Headers(ctx)
	if err == nil {
		t.Fatal("expected error for missing user_key")
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

func TestBuildConfigFingerprint(t *testing.T) {
	cfg := &config.Config{
		UserKey:      "key1",
		PluginID:     "id1",
		PluginSecret: "secret1",
		BaseURL:      "url1",
	}

	fp := auth.BuildConfigFingerprint(cfg)
	if fp == "" {
		t.Error("expected non-empty fingerprint")
	}
	if len(fp) != 64 {
		t.Errorf("expected SHA256 length 64, got %d", len(fp))
	}

	// Same config produces same fingerprint
	fp2 := auth.BuildConfigFingerprint(cfg)
	if fp != fp2 {
		t.Error("expected same fingerprint for same config")
	}

	// Different config produces different fingerprint
	cfg2 := &config.Config{
		UserKey:      "key2", // changed
		PluginID:     "id1",
		PluginSecret: "secret1",
		BaseURL:      "url1",
	}
	fp3 := auth.BuildConfigFingerprint(cfg2)
	if fp == fp3 {
		t.Error("expected different fingerprint for different config")
	}

	// nil config returns empty
	if auth.BuildConfigFingerprint(nil) != "" {
		t.Error("expected empty fingerprint for nil config")
	}
}
