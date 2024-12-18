package limiter

import (
	"context"
	"log"
	"time"
)

type RateLimiter struct {
	persistence Persistence
	ipLimit     int
	tokenLimits map[string]int
	blockTime   int
}

func NewRateLimiter(persistence Persistence, ipLimit int, blockTime int, tokenLimits map[string]int) *RateLimiter {
	return &RateLimiter{
		persistence: persistence,
		ipLimit:     ipLimit,
		blockTime:   blockTime,
		tokenLimits: tokenLimits,
	}
}

func (l *RateLimiter) Allow(ctx context.Context, key string, isToken bool) (bool, error) {
	log.Printf("RateLimiter: Verificando chave '%s', isToken: %v", key, isToken)

	isBlocked, err := l.persistence.IsBlocked(ctx, key)
	if err != nil {
		log.Printf("RateLimiter: Erro ao verificar bloqueio: %v", err)
		return false, err
	}
	if isBlocked {
		log.Printf("RateLimiter: Chave '%s' está bloqueada", key)
		return false, nil
	}

	limit := l.ipLimit
	if isToken {
		if tokenLimit, ok := l.tokenLimits[key]; ok {
			limit = tokenLimit
		}
	}

	count, err := l.persistence.Increment(ctx, key)
	if err != nil {
		log.Printf("RateLimiter: Erro ao incrementar chave '%s': %v", key, err)
		return false, err
	}
	log.Printf("RateLimiter: Contagem atual da chave '%s': %d", key, count)

	if count == 1 {
		err := l.persistence.SetExpire(ctx, key, int(time.Duration(l.blockTime).Seconds()))
		if err != nil {
			log.Printf("RateLimiter: Erro ao definir expiração para chave '%s': %v", key, err)
			return false, err
		}
	}

	if count > limit {
		log.Printf("RateLimiter: Chave '%s' atingiu o limite de %d, bloqueando", key, limit)
		_ = l.persistence.Block(ctx, key, time.Duration(l.blockTime)*time.Second)
		return false, nil
	}

	return true, nil
}
