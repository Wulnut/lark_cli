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
		Version:   session.CurrentVersion,
		LoginType: "user_key",
		UserKey:   "ou_test123",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := store.Save(ctx, s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load(ctx)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.UserKey != s.UserKey {
		t.Errorf("UserKey = %q, want %q", loaded.UserKey, s.UserKey)
	}
	if loaded.LoginType != "user_key" {
		t.Errorf("LoginType = %q", loaded.LoginType)
	}
}

func TestFileStore_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.json")
	store := session.NewFileStore(path)
	ctx := context.Background()

	s := &session.Session{Version: session.CurrentVersion, UserKey: "ou_x", CreatedAt: time.Now(), UpdatedAt: time.Now()}
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

	s := &session.Session{Version: session.CurrentVersion, UserKey: "ou_x", CreatedAt: time.Now(), UpdatedAt: time.Now()}
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

	s := &session.Session{Version: session.CurrentVersion, UserKey: "ou_x", CreatedAt: time.Now(), UpdatedAt: time.Now()}
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

	s := &session.Session{Version: session.CurrentVersion, UserKey: "ou_x", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_ = store.Save(ctx, s)

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".tmp" {
			t.Errorf("leftover tmp file: %s", e.Name())
		}
	}
}
