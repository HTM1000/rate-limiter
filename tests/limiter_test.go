package tests

import (
	"context"
	"testing"
	"time"

	"github.com/htm1000/rate-limiter/limiter"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterByIP(t *testing.T) {
	ctx := context.Background()
	redisLimiter := limiter.NewRedisLimiter("localhost", "6379")
	rateLimiter := limiter.NewRateLimiter(redisLimiter, 5, 60, nil)

	key := "test:ip"

	for i := 1; i <= 6; i++ {
		allowed, err := rateLimiter.Allow(ctx, key, false)
		assert.NoError(t, err)
		if i > 5 && allowed {
			t.Fatalf("Esperado bloquear requisição %d, mas foi permitido", i)
		}
	}

	time.Sleep(70 * time.Second)

	allowed, err := rateLimiter.Allow(ctx, key, false)
	assert.NoError(t, err)
	assert.True(t, allowed, "Esperado permitir requisição após o tempo de bloqueio")
}

func TestRateLimiterByToken(t *testing.T) {
	ctx := context.Background()
	redisLimiter := limiter.NewRedisLimiter("localhost", "6379")
	tokenLimits := map[string]int{
		"token:token1": 10,
	}
	rateLimiter := limiter.NewRateLimiter(redisLimiter, 5, 60, tokenLimits)

	tokenKey := "token:token1"

	for i := 1; i <= 11; i++ {
    allowed, err := rateLimiter.Allow(ctx, tokenKey, true)
    if err != nil {
        t.Fatal(err)
    }

    if i > 10 && allowed {
        t.Fatalf("Esperado bloquear requisição %d, mas foi permitido", i)
    }
	}

	time.Sleep(70 * time.Second)

	allowed, err := rateLimiter.Allow(ctx, tokenKey, true)
	assert.NoError(t, err)
	assert.True(t, allowed, "Esperado permitir requisição após o tempo de bloqueio")
}

func TestRateLimiterCustomTokenLimits(t *testing.T) {
	ctx := context.Background()
	redisLimiter := limiter.NewRedisLimiter("localhost", "6379")
	tokenLimits := map[string]int{
		"tokenA": 20,
		"tokenB": 5,
	}
	rateLimiter := limiter.NewRateLimiter(redisLimiter, 5, 60, tokenLimits)

	tests := []struct {
		tokenKey string
		limit    int
	}{
		{"token:tokenA", 20},
		{"token:tokenB", 5},
	}

	for _, test := range tests {
		t.Run(test.tokenKey, func(t *testing.T) {
			for i := 1; i <= test.limit+1; i++ {
				allowed, err := rateLimiter.Allow(ctx, test.tokenKey, true)
				if err != nil {
					t.Fatal(err)
				}

				if i > test.limit && allowed {
					t.Fatalf("Esperado bloquear requisição %d, mas foi permitido", i)
				}
			}
		})
	}
}

func TestRateLimiterRedisConnectionFailure(t *testing.T) {
	ctx := context.Background()
	redisLimiter := limiter.NewRedisLimiter("invalid_host", "9999")
	rateLimiter := limiter.NewRateLimiter(redisLimiter, 5, 60, nil)

	key := "test:ip"
	_, err := rateLimiter.Allow(ctx, key, false)

	assert.Error(t, err, "Esperado erro devido a falha de conexão com Redis")
}
