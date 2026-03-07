package repository

import (
	"testing"
)

func TestHashToken(t *testing.T) {
	token := "test-token-12345"
	
	hash1 := HashToken(token)
	hash2 := HashToken(token)
	
	if hash1 != hash2 {
		t.Error("same token should produce same hash")
	}
	
	if hash1 == token {
		t.Error("hash should not equal original token")
	}
	
	if len(hash1) != 64 {
		t.Errorf("SHA256 hash should be 64 characters, got %d", len(hash1))
	}
}

func TestHashToken_DifferentTokens(t *testing.T) {
	token1 := "token1"
	token2 := "token2"
	
	hash1 := HashToken(token1)
	hash2 := HashToken(token2)
	
	if hash1 == hash2 {
		t.Error("different tokens should produce different hashes")
	}
}

func TestHashToken_EmptyToken(t *testing.T) {
	hash := HashToken("")
	
	if hash == "" {
		t.Error("empty token should still produce a hash")
	}
}

func TestHashToken_Consistency(t *testing.T) {
	tokens := []string{
		"short",
		"very-long-token-with-lots-of-characters-1234567890",
		"token-with-special-chars!@#$%",
		"unicode-token",
	}
	
	for _, token := range tokens {
		hash1 := HashToken(token)
		hash2 := HashToken(token)
		
		if hash1 != hash2 {
			t.Errorf("hash should be consistent for token %q", token)
		}
	}
}
