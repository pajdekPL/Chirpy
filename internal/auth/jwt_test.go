package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tokenSecret := "12345"
	tokenSecret2 := "12345aaabbbccc"
	expiresInValid, _ := time.ParseDuration("1h")
	expiresInExpired, _ := time.ParseDuration("-1h")
	userUUID := uuid.New()

	tests := []struct {
		name                string
		userID              uuid.UUID
		expiresIn           time.Duration
		secretMakingJWT     string
		secretValidatingJWT string
		wantErr             bool
		errMsg              string
	}{
		{
			name:                "Proper JWT",
			userID:              userUUID,
			expiresIn:           expiresInValid,
			secretMakingJWT:     tokenSecret,
			secretValidatingJWT: tokenSecret,
			wantErr:             false,
			errMsg:              "",
		},
		{
			name:                "Different secrets",
			userID:              userUUID,
			expiresIn:           expiresInValid,
			secretMakingJWT:     tokenSecret,
			secretValidatingJWT: tokenSecret2,
			wantErr:             true,
			errMsg:              "signature is invalid",
		},
		{
			name:                "Expired token",
			userID:              userUUID,
			expiresIn:           expiresInExpired,
			secretMakingJWT:     tokenSecret,
			secretValidatingJWT: tokenSecret,
			wantErr:             true,
			errMsg:              "token is expired",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.secretMakingJWT, tt.expiresIn)
			if err != nil {
				t.Errorf("MakeJWT() error = %v", err)
			}

			decodedUuid, err := ValidateJWT(token, tt.secretValidatingJWT)
			if err != nil {
				if tt.wantErr {
					if !strings.Contains(err.Error(), tt.errMsg) {
						t.Errorf("ValidateJWT() expected error: '%v' to contain: '%v'", err, tt.errMsg)
					}
				}
			}
			if !tt.wantErr && decodedUuid != userUUID {
				t.Errorf("Token %v != expected: %v", decodedUuid, userUUID)
			}
		})
	}
}
