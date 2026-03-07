package config

import (
	"os"
	"testing"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	os.Clearenv()
	
	LoadConfig()
	
	if AppConfig.ServerPort != "8080" {
		t.Errorf("expected default ServerPort 8080, got %s", AppConfig.ServerPort)
	}
	if AppConfig.GinMode != "debug" {
		t.Errorf("expected default GinMode debug, got %s", AppConfig.GinMode)
	}
	if AppConfig.JWTExpireHours != 24 {
		t.Errorf("expected default JWTExpireHours 24, got %d", AppConfig.JWTExpireHours)
	}
	if AppConfig.JWTRefreshExpireDays != 7 {
		t.Errorf("expected default JWTRefreshExpireDays 7, got %d", AppConfig.JWTRefreshExpireDays)
	}
}

func TestValidate_MissingDBPassword(t *testing.T) {
	os.Clearenv()
	LoadConfig()
	
	err := Validate()
	if err == nil {
		t.Error("expected error for missing DB_PASSWORD")
	}
	if err != nil && err.Error() != "DB_PASSWORD environment variable is required" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestValidate_MissingJWTSecret(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_PASSWORD", "test_password")
	LoadConfig()
	
	err := Validate()
	if err == nil {
		t.Error("expected error for missing JWT_SECRET")
	}
	if err != nil && err.Error() != "JWT_SECRET environment variable is required" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestValidate_ShortJWTSecret(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("JWT_SECRET", "short")
	LoadConfig()
	
	err := Validate()
	if err == nil {
		t.Error("expected error for short JWT_SECRET")
	}
	if err != nil && err.Error() != "JWT_SECRET must be at least 32 characters long" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestValidate_Success(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_PASSWORD", "test_password_123")
	os.Setenv("JWT_SECRET", "this_is_a_very_long_jwt_secret_key_for_testing")
	LoadConfig()
	
	err := Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseAllowedOrigins(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"http://localhost:5173", []string{"http://localhost:5173"}},
		{"http://localhost:5173,http://localhost:3000", []string{"http://localhost:5173", "http://localhost:3000"}},
		{"", []string{}},
		{"http://example.com,https://example.com", []string{"http://example.com", "https://example.com"}},
	}
	
	for _, test := range tests {
		result := parseAllowedOrigins(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("for input %q, expected %d origins, got %d", test.input, len(test.expected), len(result))
			continue
		}
		for i, origin := range result {
			if origin != test.expected[i] {
				t.Errorf("for input %q, expected origin[%d] to be %q, got %q", test.input, i, test.expected[i], origin)
			}
		}
	}
}

func TestGetEnv(t *testing.T) {
	os.Clearenv()
	
	result := getEnv("NON_EXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("expected default_value, got %s", result)
	}
	
	os.Setenv("EXISTENT_VAR", "actual_value")
	result = getEnv("EXISTENT_VAR", "default_value")
	if result != "actual_value" {
		t.Errorf("expected actual_value, got %s", result)
	}
}
