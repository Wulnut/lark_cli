package cmdtest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"lark_cli/cmd"
	"lark_cli/internal/auth"
	"lark_cli/internal/config"
	"lark_cli/internal/session"
)

// mockConfigStore is a test double for auth.ConfigStore.
type mockConfigStore struct {
	cfg *config.Config
}

func (m *mockConfigStore) Load(ctx context.Context) (*config.Config, error) {
	return m.cfg, nil
}

func newTestDeps(t *testing.T) (cmd.Deps, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)

	var stdout, stderr bytes.Buffer
	deps := cmd.Deps{
		Config: config.Config{
			SessionPath: path,
			BaseURL:    "https://test.example.com",
			UserKey:    "",
		},
		Store:  store,
		Stdout: &stdout,
		Stderr: &stderr,
	}
	return deps, path
}

// --- Login Tests ---

func TestLogin_ValidUserKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := t.TempDir()
	sessionPath := filepath.Join(dir, "session.json")
	store := session.NewFileStore(sessionPath)
	var stdout, stderr bytes.Buffer
	deps := cmd.Deps{
		Config: config.Config{
			SessionPath: sessionPath,
			BaseURL:    "https://test.example.com",
			UserKey:    "",
		},
		Store:  store,
		Stdout: &stdout,
		Stderr: &stderr,
	}

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"login", "-w", "user_key", "ou_test123"})
	if err := root.Execute(); err != nil {
		t.Fatalf("login: %v", err)
	}

	// Verify config.json was created with user_key
	configPath := filepath.Join(home, ".lark", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	if cfg["user_key"] != "ou_test123" {
		t.Errorf("user_key = %q, want %q", cfg["user_key"], "ou_test123")
	}
}

func TestLogin_NonOuPrefixAllowed(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := t.TempDir()
	sessionPath := filepath.Join(dir, "session.json")
	store := session.NewFileStore(sessionPath)
	var stdout, stderr bytes.Buffer
	deps := cmd.Deps{
		Config: config.Config{
			SessionPath: sessionPath,
			BaseURL:    "https://test.example.com",
			UserKey:    "",
		},
		Store:  store,
		Stdout: &stdout,
		Stderr: &stderr,
	}

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"login", "-w", "user_key", "7387857889332969475"})
	if err := root.Execute(); err != nil {
		t.Fatalf("login should accept non-ou key: %v", err)
	}

	// Verify config.json was created
	configPath := filepath.Join(home, ".lark", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	if cfg["user_key"] != "7387857889332969475" {
		t.Errorf("user_key = %q", cfg["user_key"])
	}
}

func TestLogin_PreservesExistingConfigFieldsWhenUpdatingUserKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configDir := filepath.Join(home, ".lark")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	seed := map[string]any{
		"base_url":      "https://seed.example.com",
		"plugin_id":     "seed_pid",
		"plugin_secret": "seed_secret",
		"user_key":      "old_user_key",
		"session_path":  "/tmp/seed-session.json",
	}
	seedBytes, _ := json.Marshal(seed)
	if err := os.WriteFile(configPath, seedBytes, 0o600); err != nil {
		t.Fatalf("write seed config: %v", err)
	}

	deps, _ := newTestDeps(t)
	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"login", "-w", "user_key", "new_user_key"})
	if err := root.Execute(); err != nil {
		t.Fatalf("login: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read merged config: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal merged config: %v", err)
	}

	if got["user_key"] != "new_user_key" {
		t.Errorf("user_key = %v, want new_user_key", got["user_key"])
	}
	if got["base_url"] != "https://seed.example.com" {
		t.Errorf("base_url changed: %v", got["base_url"])
	}
	if got["plugin_id"] != "seed_pid" {
		t.Errorf("plugin_id changed: %v", got["plugin_id"])
	}
	if got["plugin_secret"] != "seed_secret" {
		t.Errorf("plugin_secret changed: %v", got["plugin_secret"])
	}
	if got["session_path"] != "/tmp/seed-session.json" {
		t.Errorf("session_path changed: %v", got["session_path"])
	}
}

func TestLogin_EmptyKey(t *testing.T) {
	deps, _ := newTestDeps(t)

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"login", "-w", "user_key", ""})
	root.SilenceUsage = true
	root.SilenceErrors = true
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestLogin_UnsupportedMethod(t *testing.T) {
	deps, _ := newTestDeps(t)

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"login", "-w", "oauth", "ou_test"})
	root.SilenceUsage = true
	root.SilenceErrors = true
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for unsupported method")
	}
}

// --- Logout Tests ---

func TestLogout_ExistingSession(t *testing.T) {
	deps, path := newTestDeps(t)
	store := session.NewFileStore(path)

	// Create a session first
	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = store.Save(context.Background(), sess)

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"logout"})
	if err := root.Execute(); err != nil {
		t.Fatalf("logout: %v", err)
	}

	exists, _ := store.Exists(context.Background())
	if exists {
		t.Error("session should not exist after logout")
	}
}

func TestLogout_NoSession(t *testing.T) {
	deps, _ := newTestDeps(t)

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"logout"})
	if err := root.Execute(); err != nil {
		t.Fatalf("logout should succeed when no session: %v", err)
	}
}

// --- Auth Status Tests ---

func TestAuthStatus_NotLoggedIn(t *testing.T) {
	deps, _ := newTestDeps(t)

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"auth", "status"})
	root.SilenceUsage = true
	root.SilenceErrors = true

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when not logged in")
	}

	stderr := deps.Stderr.(*bytes.Buffer)
	if !bytes.Contains(stderr.Bytes(), []byte("Not logged in")) {
		t.Errorf("stderr = %q, want 'Not logged in'", stderr.String())
	}
}

func TestAuthStatus_LoggedIn_NoToken(t *testing.T) {
	deps, path := newTestDeps(t)
	store := session.NewFileStore(path)

	deps.Config.UserKey = "ou_statustest"

	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = store.Save(context.Background(), sess)

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"auth", "status"})
	if err := root.Execute(); err != nil {
		t.Fatalf("auth status: %v", err)
	}

	stdout := deps.Stdout.(*bytes.Buffer)
	output := stdout.String()
	if !bytes.Contains([]byte(output), []byte("ou_sta***")) {
		t.Errorf("expected masked user_key, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Plugin token: none")) {
		t.Errorf("expected 'Plugin token: none', got: %s", output)
	}
}

func TestAuthStatus_WithRefresh(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"data":  map[string]any{"token": "p-refreshed-token-123", "expire_time": 7200},
			"error": map[string]any{"code": 0, "msg": "success"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)

	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = store.Save(context.Background(), sess)

	cfgStore := &mockConfigStore{
		cfg: &config.Config{
			UserKey:      "ou_refresh_test",
			PluginID:     "pid",
			PluginSecret: "psecret",
			BaseURL:      server.URL,
		},
	}

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	tokenProvider := auth.NewPluginTokenProvider(cfgStore, store, client, 10*time.Minute)
	headerProvider := auth.NewHeaderProvider(tokenProvider)

	var stdout, stderr bytes.Buffer
	deps := cmd.Deps{
		Config: config.Config{
			SessionPath: path,
			BaseURL:     server.URL,
			UserKey:     "ou_refresh_test",
		},
		Store:               store,
		PluginTokenProvider: tokenProvider,
		HeaderProvider:      headerProvider,
		Stdout:              &stdout,
		Stderr:              &stderr,
	}

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"auth", "status", "--refresh-plugin-token"})
	if err := root.Execute(); err != nil {
		t.Fatalf("auth status --refresh: %v", err)
	}

	output := stdout.String()
	if !bytes.Contains([]byte(output), []byte("refreshed")) {
		t.Errorf("expected 'refreshed' in output, got: %s", output)
	}

	// Verify token was persisted
	loaded, _ := store.Load(context.Background())
	if loaded.PluginAccessToken != "p-refreshed-token-123" {
		t.Errorf("persisted token = %q", loaded.PluginAccessToken)
	}
}
