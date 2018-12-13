package app

import "os"

// GetEnv retrieves the value of the environment variable named by the key
// or fallback if not found
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
