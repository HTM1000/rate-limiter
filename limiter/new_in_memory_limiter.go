package limiter

import (
	"context"
	"sync"
	"time"
)

type InMemoryLimiter struct {
	data      map[string]int
	blocked   map[string]time.Time
	mu        sync.Mutex
}

func NewInMemoryLimiter() *InMemoryLimiter {
	return &InMemoryLimiter{
		data:    make(map[string]int),
		blocked: make(map[string]time.Time),
	}
}

func (l *InMemoryLimiter) Increment(ctx context.Context, key string) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.data[key]; !exists {
		l.data[key] = 0
	}
	l.data[key]++
	return l.data[key], nil
}

func (l *InMemoryLimiter) Get(ctx context.Context, key string) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	count, exists := l.data[key]
	if !exists {
		return 0, nil 
	}
	return count, nil
}

func (l *InMemoryLimiter) Reset(ctx context.Context, key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.data, key)
	return nil
}

func (l *InMemoryLimiter) TTL(ctx context.Context, key string) (int, error) {
	return 60, nil
}

func (l *InMemoryLimiter) SetExpire(ctx context.Context, key string, seconds int) error {
	return nil
}

func (l *InMemoryLimiter) Block(ctx context.Context, key string, duration time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.blocked[key] = time.Now().Add(duration)
	return nil
}

func (l *InMemoryLimiter) IsBlocked(ctx context.Context, key string) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	blockTime, exists := l.blocked[key]
	if !exists {
		return false, nil
	}

	if time.Now().After(blockTime) {
		delete(l.blocked, key)
		return false, nil
	}

	return true, nil
}
