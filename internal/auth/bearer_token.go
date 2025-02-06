package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("bearer token not found")
	}
	elements := strings.Split(bearer, " ")

	if len(elements) != 2 {
		return "", fmt.Errorf("bearer token not found")
	}
	if strings.ToLower(elements[0]) != "bearer" {
		return "", fmt.Errorf("bearer token not found")
	}

	return elements[1], nil
}
