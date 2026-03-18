package auth

import (
	"context"
	"errors"
	"time"

	"lark_cli/internal/session"
)

var ErrNotLoggedIn = errors.New("not logged in: run 'lark login' first")

// Clock abstracts time for testing.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// PluginTokenProvider manages cached plugin_access_token.
type PluginTokenProvider interface {
	Get(ctx context.Context) (string, error)
	ForceRefresh(ctx context.Context) (string, error)
}

type pluginTokenProvider struct {
	store  session.Store
	client *PluginTokenClient
	clock  Clock
	leeway time.Duration
}

// NewPluginTokenProvider creates a provider that caches plugin tokens via the session store.
func NewPluginTokenProvider(store session.Store, client *PluginTokenClient, leeway time.Duration) PluginTokenProvider {
	return &pluginTokenProvider{
		store:  store,
		client: client,
		clock:  realClock{},
		leeway: leeway,
	}
}

// NewPluginTokenProviderWithClock is used for testing with a fake clock.
func NewPluginTokenProviderWithClock(store session.Store, client *PluginTokenClient, leeway time.Duration, clock Clock) PluginTokenProvider {
	return &pluginTokenProvider{
		store:  store,
		client: client,
		clock:  clock,
		leeway: leeway,
	}
}

func (p *pluginTokenProvider) Get(ctx context.Context) (string, error) {
	sess, err := p.store.Load(ctx)
	if err != nil {
		return "", ErrNotLoggedIn
	}
	if sess.UserKey == "" {
		return "", ErrNotLoggedIn
	}

	if sess.PluginAccessToken != "" && p.clock.Now().Before(sess.PluginAccessTokenExpiresAt.Add(-p.leeway)) {
		return sess.PluginAccessToken, nil
	}

	return p.refresh(ctx, sess)
}

func (p *pluginTokenProvider) ForceRefresh(ctx context.Context) (string, error) {
	sess, err := p.store.Load(ctx)
	if err != nil {
		return "", ErrNotLoggedIn
	}
	if sess.UserKey == "" {
		return "", ErrNotLoggedIn
	}
	return p.refresh(ctx, sess)
}

func (p *pluginTokenProvider) refresh(ctx context.Context, sess *session.Session) (string, error) {
	token, expiresIn, err := p.client.Fetch(ctx)
	if err != nil {
		return "", err
	}

	sess.PluginAccessToken = token
	sess.PluginAccessTokenExpiresAt = p.clock.Now().Add(expiresIn)
	sess.UpdatedAt = p.clock.Now()

	if err := p.store.Save(ctx, sess); err != nil {
		return "", err
	}

	return token, nil
}
