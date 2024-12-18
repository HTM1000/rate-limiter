package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Aviso: Não foi possível converter %s para int, usando valor padrão %d\n", key, defaultValue)
		return defaultValue
	}

	return value
}

func ParseEnvTokenLimits(key string) map[string]int {
	tokenLimits := make(map[string]int)

	valueStr := os.Getenv(key)
	if valueStr == "" {
		return tokenLimits
	}

	for _, tokenConfig := range strings.Split(valueStr, ",") {
		parts := strings.Split(tokenConfig, ":")
		if len(parts) == 2 {
			if limit, err := strconv.Atoi(parts[1]); err == nil {
				tokenLimits[parts[0]] = limit
			} else {
				log.Printf("Aviso: Não foi possível converter o limite de %s, ignorando", tokenConfig)
			}
		}
	}

	return tokenLimits
}
