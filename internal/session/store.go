package session

import "context"

type Store interface {
	Load(ctx context.Context) (*Session, error)
	Save(ctx context.Context, s *Session) error
	Delete(ctx context.Context) error
	Exists(ctx context.Context) (bool, error)
}
