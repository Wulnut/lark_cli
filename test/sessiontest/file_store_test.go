package sessiontest

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"lark_cli/internal/session"
)

func TestFileStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	now := time.Now().Truncate(time.Second)
	s := &session.Session{
		Version:              session.CurrentVersion,
		LoginType:           "user_key",
		ConfigFingerprint:    "fp123",
		PluginAccessToken:   "token-abc",
		PluginAccessTokenExpiresAt: now.Add(2 * time.Hour),
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := store.Save(ctx, s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load(ctx)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.LoginType != "user_key" {
		t.Errorf("LoginType = %q", loaded.LoginType)
	}
	if loaded.ConfigFingerprint != "fp123" {
		t.Errorf("ConfigFingerprint = %q", loaded.ConfigFingerprint)
	}
	if loaded.PluginAccessToken != "token-abc" {
		t.Errorf("PluginAccessToken = %q", loaded.PluginAccessToken)
	}
}

func TestFileStore_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	s := &session.Session{Version: session.CurrentVersion, LoginType: "user_key", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := store.Save(ctx, s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	perm := info.Mode().Perm()
	if perm != fs.FileMode(0o600) {
		t.Errorf("file perm = %o, want 0600", perm)
	}
}

func TestFileStore_Delete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	s := &session.Session{Version: session.CurrentVersion, LoginType: "user_key", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_ = store.Save(ctx, s)

	if err := store.Delete(ctx); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	exists, err := store.Exists(ctx)
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if exists {
		t.Error("session should not exist after delete")
	}
}

func TestFileStore_DeleteNonExistent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	if err := store.Delete(ctx); err != nil {
		t.Fatalf("Delete nonexistent should not error: %v", err)
	}
}

func TestFileStore_Exists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	exists, _ := store.Exists(ctx)
	if exists {
		t.Error("should not exist initially")
	}

	s := &session.Session{Version: session.CurrentVersion, LoginType: "user_key", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_ = store.Save(ctx, s)

	exists, _ = store.Exists(ctx)
	if !exists {
		t.Error("should exist after save")
	}
}

func TestFileStore_NoLeftoverTmpFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	s := &session.Session{Version: session.CurrentVersion, LoginType: "user_key", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_ = store.Save(ctx, s)

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".tmp" {
			t.Errorf("leftover tmp file: %s", e.Name())
		}
	}
}

func TestSession_IsValid(t *testing.T) {
	now := int64(1000000)
	validFingerprint := "fp123"

	tests := []struct {
		name        string
		sess        *session.Session
		now         int64
		fingerprint string
		want        bool
	}{
		{
			name: "valid session",
			sess: &session.Session{
				PluginAccessToken:          "tok",
				PluginAccessTokenExpiresAt: time.Unix(now+3600, 0),
				ConfigFingerprint:           validFingerprint,
			},
			now:         now,
			fingerprint: validFingerprint,
			want:        true,
		},
		{
			name:        "nil session",
			sess:        nil,
			now:         now,
			fingerprint: validFingerprint,
			want:        false,
		},
		{
			name: "empty token",
			sess: &session.Session{
				PluginAccessToken:          "",
				PluginAccessTokenExpiresAt: time.Unix(now+3600, 0),
				ConfigFingerprint:          validFingerprint,
			},
			now:         now,
			fingerprint: validFingerprint,
			want:        false,
		},
		{
			name: "expired token",
			sess: &session.Session{
				PluginAccessToken:          "tok",
				PluginAccessTokenExpiresAt: time.Unix(now-1, 0),
				ConfigFingerprint:          validFingerprint,
			},
			now:         now,
			fingerprint: validFingerprint,
			want:        false,
		},
		{
			name: "fingerprint mismatch",
			sess: &session.Session{
				PluginAccessToken:          "tok",
				PluginAccessTokenExpiresAt: time.Unix(now+3600, 0),
				ConfigFingerprint:          "old_fp",
			},
			now:         now,
			fingerprint: validFingerprint,
			want:        false,
		},
		{
			name: "exactly at expiry",
			sess: &session.Session{
				PluginAccessToken:          "tok",
				PluginAccessTokenExpiresAt: time.Unix(now, 0),
				ConfigFingerprint:          validFingerprint,
			},
			now:         now,
			fingerprint: validFingerprint,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sess.IsValid(tt.now, tt.fingerprint)
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		sess *session.Session
		want bool
	}{
		{"nil session", nil, true},
		{"empty token", &session.Session{}, true},
		{"has token", &session.Session{PluginAccessToken: "tok"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sess.IsEmpty()
			if got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
