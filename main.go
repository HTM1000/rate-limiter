package main

import (
	"log"
	"net/http"
	"os"

	"github.com/htm1000/rate-limiter/limiter"
	"github.com/htm1000/rate-limiter/middleware"
	"github.com/htm1000/rate-limiter/utils"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	redisHost := os.Getenv("REDIS_HOST")

	log.Printf(redisHost)

	redisPort := os.Getenv("REDIS_PORT")

	log.Printf(redisPort)

	ipLimit := utils.GetEnvAsInt("REQUEST_LIMIT_IP", 5)
	blockTime := utils.GetEnvAsInt("BLOCK_TIME", 300)
	tokenLimits := utils.ParseEnvTokenLimits("TOKEN_LIMITS")

	log.Printf("Conectando ao Redis em %s:%s...", redisHost, redisPort)

	redisLimiter := limiter.NewRedisLimiter(redisHost, redisPort)
	if redisLimiter == nil {
		log.Fatalf("Erro ao conectar ao Redis em %s:%s", redisHost, redisPort)
	}

	rateLimiter := limiter.NewRateLimiter(redisLimiter, ipLimit, blockTime, tokenLimits)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Recebida requisição em /")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request OK"))
	})

	log.Println("Servidor iniciado em http://localhost:8081")
	err = http.ListenAndServe(":8081", middleware.RateLimiterMiddleware(rateLimiter)(mux))
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
