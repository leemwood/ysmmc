package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	
	if hash == "" {
		t.Error("hash should not be empty")
	}
	if hash == password {
		t.Error("hash should not equal plain password")
	}
	if len(hash) < 20 {
		t.Error("hash should be sufficiently long")
	}
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	password := "testPassword123"
	
	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)
	
	if hash1 == hash2 {
		t.Error("different hashes should be generated for the same password (bcrypt salt)")
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	password := "testPassword123"
	hash, _ := HashPassword(password)
	
	if !CheckPassword(password, hash) {
		t.Error("password check should succeed for correct password")
	}
}

func TestCheckPassword_Incorrect(t *testing.T) {
	password := "testPassword123"
	wrongPassword := "wrongPassword456"
	hash, _ := HashPassword(password)
	
	if CheckPassword(wrongPassword, hash) {
		t.Error("password check should fail for incorrect password")
	}
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	hash, _ := HashPassword("somePassword")
	
	if CheckPassword("", hash) {
		t.Error("empty password should not match")
	}
}

func TestCheckPassword_EmptyHash(t *testing.T) {
	if CheckPassword("somePassword", "") {
		t.Error("should return false for empty hash")
	}
}

func TestHashPassword_VariousPasswords(t *testing.T) {
	passwords := []string{
		"short",
		"verylongpasswordwithlotsofcharacters1234567890",
		"password with spaces",
		"special!@#$%^&*()chars",
		"unicode密码",
	}
	
	for _, password := range passwords {
		hash, err := HashPassword(password)
		if err != nil {
			t.Errorf("failed to hash password %q: %v", password, err)
			continue
		}
		
		if !CheckPassword(password, hash) {
			t.Errorf("password check failed for %q", password)
		}
	}
}
