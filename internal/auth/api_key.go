package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("ApiKey not found")
	}
	elements := strings.Split(apiKey, " ")

	if len(elements) != 2 {
		return "", fmt.Errorf("ApiKey not found")
	}
	if strings.ToLower(elements[0]) != "apikey" {
		return "", fmt.Errorf("ApiKey not found")
	}

	return elements[1], nil
}
