package limiter

import (
	"context"
	"fmt"
	"time"
	"strconv" 

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
	if err != nil {
		return 0, err
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if ttl <= 0 {
		err := r.client.Expire(ctx, key, time.Second*60).Err()
		if err != nil {
			return 0, err
		}
	}

	return int(result), nil
}

func (r *RedisLimiter) Get(ctx context.Context, key string) (int, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {	
		return 0, err
	}

	return strconv.Atoi(value)
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
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
			value = "0" 
	} else if err != nil {
			return err
	}
	return r.client.Set(ctx, key, value, time.Duration(seconds)*time.Second).Err()
}
