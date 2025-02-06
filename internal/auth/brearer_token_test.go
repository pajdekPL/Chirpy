package auth

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedToken string
		err           error
	}{
		{
			name: "Proper Bearer Token",
			headers: http.Header{
				"Authorization": []string{"Bearer some_token"},
				"Content-Type":  []string{"application/json"},
			},
			expectedToken: "some_token",
			err:           nil,
		},
		{
			name:          "No authorization header",
			headers:       http.Header{},
			expectedToken: "",
			err:           fmt.Errorf("bearer token not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			if !errors.Is(err, tt.err) {
				t.Errorf("GetBearerToken() expected error: %v != %v", tt.err, err)
			}

			if token != tt.expectedToken {
				t.Errorf("GetBearerToken() expected token: %v != %v", tt.expectedToken, token)

			}

		})
	}
}
