package auth

import (
	"context"

	"lark_cli/internal/session"
)

// HeaderProvider assembles auth headers for Feishu Project OpenAPI calls.
type HeaderProvider interface {
	Headers(ctx context.Context) (map[string]string, error)
}

type headerProvider struct {
	store         session.Store
	tokenProvider PluginTokenProvider
}

// NewHeaderProvider creates a provider that returns X-Plugin-Token and X-User-Key headers.
func NewHeaderProvider(store session.Store, tokenProvider PluginTokenProvider) HeaderProvider {
	return &headerProvider{
		store:         store,
		tokenProvider: tokenProvider,
	}
}

func (h *headerProvider) Headers(ctx context.Context) (map[string]string, error) {
	sess, err := h.store.Load(ctx)
	if err != nil {
		return nil, ErrNotLoggedIn
	}
	if sess.UserKey == "" {
		return nil, ErrNotLoggedIn
	}

	token, err := h.tokenProvider.Get(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"X-Plugin-Token": token,
		"X-User-Key":     sess.UserKey,
	}, nil
}
