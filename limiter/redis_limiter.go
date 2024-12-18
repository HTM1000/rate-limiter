package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(host string, port string) *RedisLimiter {
	return &RedisLimiter{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", host, port),
		}),
	}
}

func (r *RedisLimiter) Increment(ctx context.Context, key string) (int, error) {
	result, err := r.client.Incr(ctx, key).Result()
	return int(result), err
}

func (r *RedisLimiter) Reset(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisLimiter) IsBlocked(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key+":blocked").Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisLimiter) Block(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Set(ctx, key+":blocked", true, duration).Err()
}

func (r *RedisLimiter) TTL(ctx context.Context, key string) (int, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(ttl.Seconds()), nil
}

func (r *RedisLimiter) SetExpire(ctx context.Context, key string, seconds int) error {
	return r.client.Expire(ctx, key, time.Duration(seconds)*time.Second).Err()
}
