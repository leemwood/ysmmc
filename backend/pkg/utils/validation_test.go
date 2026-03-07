package utils

import (
	"testing"
)

func TestSanitizeString_HTML(t *testing.T) {
	input := "<script>alert('xss')</script>"
	result := SanitizeString(input)
	
	if result == input {
		t.Error("HTML should be escaped")
	}
	if contains(result, "<script>") {
		t.Error("script tags should be removed or escaped")
	}
}

func TestSanitizeString_Script(t *testing.T) {
	input := `<script>document.cookie</script>`
	result := SanitizeString(input)
	
	if contains(result, "<script>") {
		t.Error("script tags should be removed")
	}
}

func TestSanitizeString_EventHandlers(t *testing.T) {
	input := `<img src="x" onerror="alert('xss')">`
	result := SanitizeString(input)
	
	if contains(result, `onerror="alert`) {
		t.Error("event handlers should be removed")
	}
}

func TestSanitizeString_PlainText(t *testing.T) {
	input := "Hello World"
	result := SanitizeString(input)
	
	if result != input {
		t.Errorf("plain text should not be modified, got %q", result)
	}
}

func TestSanitizeString_TrimSpace(t *testing.T) {
	input := "  hello world  "
	result := SanitizeString(input)
	
	if result != "hello world" {
		t.Errorf("spaces should be trimmed, got %q", result)
	}
}

func TestStripHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<p>Hello</p>", "Hello"},
		{"<div><span>Test</span></div>", "Test"},
		{"No HTML here", "No HTML here"},
		{"<a href='link'>Link</a>", "Link"},
	}
	
	for _, test := range tests {
		result := StripHTML(test.input)
		if result != test.expected {
			t.Errorf("for input %q, expected %q, got %q", test.input, test.expected, result)
		}
	}
}

func TestValidateEmail_Valid(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@example.com",
		"user+tag@example.org",
		"user123@test.co.uk",
	}
	
	for _, email := range validEmails {
		if !ValidateEmail(email) {
			t.Errorf("email %q should be valid", email)
		}
	}
}

func TestValidateEmail_Invalid(t *testing.T) {
	invalidEmails := []string{
		"",
		"invalid",
		"invalid@",
		"@example.com",
		"user@.com",
		"user@example",
		"user @example.com",
	}
	
	for _, email := range invalidEmails {
		if ValidateEmail(email) {
			t.Errorf("email %q should be invalid", email)
		}
	}
}

func TestValidateUsername_Valid(t *testing.T) {
	validUsernames := []string{
		"user123",
		"test_user",
		"test-user",
		"TestUser",
		"ab",
		"testuser123456789012345678901234567890",
		"testuser",
		"test",
	}
	
	for _, username := range validUsernames {
		if !ValidateUsername(username) {
			t.Errorf("username %q should be valid", username)
		}
	}
}

func TestValidateUsername_Chinese(t *testing.T) {
	chineseUsernames := []string{
		"test",
		"user",
		"admin",
	}
	
	for _, username := range chineseUsernames {
		if !ValidateUsername(username) {
			t.Errorf("username %q should be valid", username)
		}
	}
}

func TestValidateUsername_Invalid(t *testing.T) {
	invalidUsernames := []string{
		"",
		"a",
		"test user",
		"test@user",
		"test!user",
		"test.user",
	}
	
	for _, username := range invalidUsernames {
		if ValidateUsername(username) {
			t.Errorf("username %q should be invalid", username)
		}
	}
}

func TestValidateUsername_TooLong(t *testing.T) {
	longUsername := "a" + string(make([]byte, 50))
	if len(longUsername) <= 50 {
		longUsername = ""
		for i := 0; i < 51; i++ {
			longUsername += "a"
		}
	}
	
	if ValidateUsername(longUsername) {
		t.Errorf("username longer than 50 characters should be invalid")
	}
}

func TestValidatePassword(t *testing.T) {
	validPasswords := []string{
		"123456",
		"password",
		"verylongpassword123",
	}
	
	for _, password := range validPasswords {
		if !ValidatePassword(password) {
			t.Errorf("password %q should be valid", password)
		}
	}
	
	invalidPasswords := []string{
		"",
		"12345",
		"short",
	}
	
	for _, password := range invalidPasswords {
		if ValidatePassword(password) {
			t.Errorf("password %q should be invalid", password)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hello"},
		{"test", 4, "test"},
		{"", 10, ""},
	}
	
	for _, test := range tests {
		result := TruncateString(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("for input %q with maxLen %d, expected %q, got %q", 
				test.input, test.maxLen, test.expected, result)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
