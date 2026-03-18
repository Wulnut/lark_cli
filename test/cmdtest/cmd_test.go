package cmdtest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"lark_cli/cmd"
	"lark_cli/internal/auth"
	"lark_cli/internal/config"
	"lark_cli/internal/session"
)

func newTestDeps(t *testing.T) (cmd.Deps, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)

	var stdout, stderr bytes.Buffer
	deps := cmd.Deps{
		Config: config.Config{
			SessionPath: path,
			BaseURL:     "https://test.example.com",
		},
		Store:  store,
		Stdout: &stdout,
		Stderr: &stderr,
	}
	return deps, path
}

// --- Login Tests ---

func TestLogin_ValidUserKey(t *testing.T) {
	deps, path := newTestDeps(t)

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"login", "-w", "user_key", "ou_test123"})
	if err := root.Execute(); err != nil {
		t.Fatalf("login: %v", err)
	}

	// Verify session was created
	store := session.NewFileStore(path)
	sess, err := store.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if sess.UserKey != "ou_test123" {
		t.Errorf("UserKey = %q", sess.UserKey)
	}
	if sess.LoginType != "user_key" {
		t.Errorf("LoginType = %q", sess.LoginType)
	}
}

func TestLogin_NonOuPrefixAllowed(t *testing.T) {
	deps, path := newTestDeps(t)

	root := cmd.NewRootCmd(deps)
	root.SetArgs([]string{"login", "-w", "user_key", "7387857889332969475"})
	if err := root.Execute(); err != nil {
		t.Fatalf("login should accept non-ou key: %v", err)
	}

	store := session.NewFileStore(path)
	sess, err := store.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if sess.UserKey != "7387857889332969475" {
		t.Errorf("UserKey = %q", sess.UserKey)
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
		UserKey:   "ou_test",
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

	sess := &session.Session{
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		UserKey:   "ou_statustest",
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
		UserKey:   "ou_refresh_test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = store.Save(context.Background(), sess)

	client := auth.NewPluginTokenClient(server.Client(), server.URL, "pid", "psecret")
	tokenProvider := auth.NewPluginTokenProvider(store, client, 10*time.Minute)
	headerProvider := auth.NewHeaderProvider(store, tokenProvider)

	var stdout, stderr bytes.Buffer
	deps := cmd.Deps{
		Config:              config.Config{SessionPath: path},
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
