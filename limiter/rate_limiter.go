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
		} else {
			log.Printf("RateLimiter: Token '%s' não encontrado", key)
		}
	}
	
	currentCount, err := l.persistence.Get(ctx, key)
	if err != nil {
		log.Printf("RateLimiter: Erro ao obter valor da chave '%s': %v", key, err)
		return false, err
	}
	log.Printf("RateLimiter: Valor atual da chave '%s': %d", key, currentCount)
	
	if currentCount >= limit {
		log.Printf("RateLimiter: Chave '%s' atingiu o limite de %d, bloqueando", key, limit)
		
		err = l.persistence.Block(ctx, key, time.Duration(l.blockTime)*time.Second)
		if err != nil {
			log.Printf("RateLimiter: Erro ao bloquear chave '%s': %v", key, err)
			return false, err
		}
		_ = l.persistence.Reset(ctx, key) 
		return false, nil
	}

	count, err := l.persistence.Increment(ctx, key)
	if err != nil {
		log.Printf("RateLimiter: Erro ao incrementar chave '%s': %v", key, err)
		return false, err
	}
	log.Printf("RateLimiter: Valor incrementado da chave '%s': %d", key, count)

	return true, nil
}
