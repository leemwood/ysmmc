package utils

import (
	"html"
	"regexp"
	"strings"
)

var (
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
	scriptRegex  = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	eventRegex   = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`)
)

func SanitizeString(input string) string {
	input = html.EscapeString(input)
	input = scriptRegex.ReplaceAllString(input, "")
	input = eventRegex.ReplaceAllString(input, "")
	return strings.TrimSpace(input)
}

func SanitizeHTML(input string) string {
	input = scriptRegex.ReplaceAllString(input, "")
	input = eventRegex.ReplaceAllString(input, "")
	return strings.TrimSpace(input)
}

func StripHTML(input string) string {
	return htmlTagRegex.ReplaceAllString(input, "")
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateUsername(username string) bool {
	if len(username) < 2 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-` + "\u4e00-\u9fa5" + `]+$`)
	return usernameRegex.MatchString(username)
}

func ValidatePassword(password string) bool {
	return len(password) >= 6
}

func TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	return input[:maxLength]
}
