package auth

import (
	"context"
	"errors"
	"time"

	"lark_cli/internal/config"
	"lark_cli/internal/session"
)

var ErrNotLoggedIn = errors.New("not logged in: run 'lark login' first")

// Clock abstracts time for testing.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// ConfigStore provides access to the current config.
type ConfigStore interface {
	Load(ctx context.Context) (*config.Config, error)
}

// AuthContext contains the authentication context needed for API requests.
type AuthContext struct {
	UserKey     string
	PluginToken string
}

// PluginTokenProvider manages cached plugin_access_token.
type PluginTokenProvider interface {
	GetAuthContext(ctx context.Context) (*AuthContext, error)
	ForceRefresh(ctx context.Context) (*AuthContext, error)
}

type pluginTokenProvider struct {
	configStore ConfigStore
	store       session.Store
	client      *PluginTokenClient
	clock       Clock
	leeway      time.Duration
}

// NewPluginTokenProvider creates a provider that caches plugin tokens via the session store.
func NewPluginTokenProvider(configStore ConfigStore, store session.Store, client *PluginTokenClient, leeway time.Duration) PluginTokenProvider {
	return &pluginTokenProvider{
		configStore: configStore,
		store:       store,
		client:      client,
		clock:       realClock{},
		leeway:      leeway,
	}
}

// NewPluginTokenProviderWithClock is used for testing with a fake clock.
func NewPluginTokenProviderWithClock(configStore ConfigStore, store session.Store, client *PluginTokenClient, leeway time.Duration, clock Clock) PluginTokenProvider {
	return &pluginTokenProvider{
		configStore: configStore,
		store:       store,
		client:      client,
		clock:       clock,
		leeway:      leeway,
	}
}

// GetAuthContext returns a valid auth context, refreshing the token if needed.
func (p *pluginTokenProvider) GetAuthContext(ctx context.Context) (*AuthContext, error) {
	cfg, err := p.configStore.Load(ctx)
	if err != nil {
		return nil, err
	}
	if err := cfg.ValidateForOpenAPI(); err != nil {
		return nil, err
	}

	fingerprint := BuildConfigFingerprint(cfg)
	now := p.clock.Now().Unix()

	sess, err := p.store.Load(ctx)
	// If session doesn't exist or is invalid, refresh will create a new one
	if err == nil && sess != nil && sess.IsValid(now, fingerprint) {
		return &AuthContext{
			UserKey:     cfg.UserKey,
			PluginToken: sess.PluginAccessToken,
		}, nil
	}

	return p.refreshLocked(ctx, cfg, fingerprint)
}

// ForceRefresh forces a token refresh regardless of cache state.
func (p *pluginTokenProvider) ForceRefresh(ctx context.Context) (*AuthContext, error) {
	cfg, err := p.configStore.Load(ctx)
	if err != nil {
		return nil, err
	}
	if err := cfg.ValidateForOpenAPI(); err != nil {
		return nil, err
	}

	fingerprint := BuildConfigFingerprint(cfg)
	return p.refreshLocked(ctx, cfg, fingerprint)
}

func (p *pluginTokenProvider) refreshLocked(ctx context.Context, cfg *config.Config, fingerprint string) (*AuthContext, error) {
	token, expiresIn, err := p.client.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	now := p.clock.Now()
	sess, err := p.store.Load(ctx)
	if err != nil {
		// If session doesn't exist, create a new one
		sess = &session.Session{
			Version:   session.CurrentVersion,
			LoginType: "user_key",
		}
	}

	sess.PluginAccessToken = token
	sess.PluginAccessTokenExpiresAt = now.Add(expiresIn)
	sess.ConfigFingerprint = fingerprint
	sess.UpdatedAt = now

	if err := p.store.Save(ctx, sess); err != nil {
		return nil, err
	}

	return &AuthContext{
		UserKey:     cfg.UserKey,
		PluginToken: token,
	}, nil
}
