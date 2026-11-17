package nexusmc

import (
	"net/url"
	"os"
	"testing"
)

func TestAuthorizeScopes(t *testing.T) {
	t.Setenv("NEXUSMC_OAUTH_SCOPES", "placeholder")
	if err := os.Unsetenv("NEXUSMC_OAUTH_SCOPES"); err != nil {
		t.Fatal(err)
	}
	u, _ := url.Parse(AuthorizeURL("https://example.com/callback", "state"))
	if u.Query().Get("scope") != defaultScopes {
		t.Fatal("missing defaults")
	}
	for _, scope := range []string{"", "   ", "user:basic user:email"} {
		t.Setenv("NEXUSMC_OAUTH_SCOPES", scope)
		u, _ := url.Parse(AuthorizeURL("https://example.com/callback", "state"))
		if scope == "" || scope == "   " {
			if u.Query().Has("scope") {
				t.Fatal("explicit empty scope must be omitted")
			}
		} else if u.Query().Get("scope") != scope {
			t.Fatal("custom scope lost")
		}
	}
}
