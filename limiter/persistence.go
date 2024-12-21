package limiter

import (
	"context"
	"time"
)

type Persistence interface {
	Increment(ctx context.Context, key string) (int, error)
	Reset(ctx context.Context, key string) error
	TTL(ctx context.Context, key string) (int, error)
	SetExpire(ctx context.Context, key string, seconds int) error
	Block(ctx context.Context, key string, duration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) (int, error)
}
