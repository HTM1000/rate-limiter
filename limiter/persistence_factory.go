// persistence_factory.go
package limiter

import (
	"fmt"
	"os"
)

func NewPersistence() (Persistence, error) {
	storage := os.Getenv("STORAGE_BACKEND")

	switch storage {
		case "redis":
			host := os.Getenv("REDIS_HOST")
			port := os.Getenv("REDIS_PORT")
			return NewRedisLimiter(host, port), nil
		case "in_memory":
			return NewInMemoryLimiter(), nil
		default:
			return nil, fmt.Errorf("storage backend '%s' não suportado", storage)
	}
}
