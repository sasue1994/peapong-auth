package security

import (
	"fmt"
	"os"
	"strings"
)

func valiationHeader(token string) bool {
	if token == "" {
		return false
	}
	return true
}

func ValidationAPIKey(apiKey string) bool {
	expectedKey := os.Getenv("API_KEY")

	fmt.Println(expectedKey)

	fmt.Println(apiKey)

	if expectedKey == "" || apiKey == "" {
		return false
	}

	return strings.TrimSpace(apiKey) == strings.TrimSpace(expectedKey)
}
