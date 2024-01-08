package main

import (
	"fmt"
	"os"
)

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("env key %s isn't defined", key)
}
