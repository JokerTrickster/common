package redis

import "os"

// getEnvOrFallback fetches an environment variable or returns a fallback value
func getEnvOrFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
