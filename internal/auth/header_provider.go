package auth

import (
	"context"
)

// HeaderProvider assembles auth headers for Feishu Project OpenAPI calls.
type HeaderProvider interface {
	Headers(ctx context.Context) (map[string]string, error)
}

type headerProvider struct {
	tokenProvider PluginTokenProvider
}

// NewHeaderProvider creates a provider that returns X-Plugin-Token and X-User-Key headers.
func NewHeaderProvider(tokenProvider PluginTokenProvider) HeaderProvider {
	return &headerProvider{
		tokenProvider: tokenProvider,
	}
}

func (h *headerProvider) Headers(ctx context.Context) (map[string]string, error) {
	authCtx, err := h.tokenProvider.GetAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"X-Plugin-Token": authCtx.PluginToken,
		"X-User-Key":     authCtx.UserKey,
	}, nil
}
