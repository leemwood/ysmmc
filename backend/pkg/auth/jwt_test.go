package auth

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
)

func init() {
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("JWT_SECRET", "test_jwt_secret_key_for_testing_purposes_32")
	config.LoadConfig()
}

func TestGenerateToken(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	
	tokens, err := GenerateToken(userID, email, role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	
	if tokens.AccessToken == "" {
		t.Error("access token should not be empty")
	}
	if tokens.RefreshToken == "" {
		t.Error("refresh token should not be empty")
	}
	if tokens.ExpiresIn <= 0 {
		t.Error("expires_in should be positive")
	}
}

func TestParseToken_Valid(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	role := "admin"
	
	tokens, err := GenerateToken(userID, email, role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	
	claims, err := ParseToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to parse valid token: %v", err)
	}
	
	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	invalidTokens := []string{
		"",
		"invalid.token.format",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
	}
	
	for _, token := range invalidTokens {
		_, err := ParseToken(token)
		if err == nil {
			t.Errorf("expected error for invalid token: %s", token)
		}
	}
}

func TestParseRefreshToken(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	
	tokens, err := GenerateToken(userID, email, role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	
	parsedUserID, err := ParseRefreshToken(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("failed to parse refresh token: %v", err)
	}
	
	if parsedUserID != userID {
		t.Errorf("expected userID %s, got %s", userID, parsedUserID)
	}
}

func TestParseRefreshToken_Invalid(t *testing.T) {
	invalidTokens := []string{
		"",
		"invalid.refresh.token",
	}
	
	for _, token := range invalidTokens {
		_, err := ParseRefreshToken(token)
		if err == nil {
			t.Errorf("expected error for invalid refresh token: %s", token)
		}
	}
}

func TestTokenExpiration(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	
	tokens, err := GenerateToken(userID, email, role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	
	claims, err := ParseToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	
	now := time.Now()
	if claims.ExpiresAt.Time.Before(now) {
		t.Error("token should not be expired immediately after generation")
	}
	
	expectedExpiry := now.Add(time.Duration(24) * time.Hour)
	tolerance := time.Minute
	
	if claims.ExpiresAt.Time.Sub(expectedExpiry) > tolerance {
		t.Errorf("token expiry time is not within expected range")
	}
}
