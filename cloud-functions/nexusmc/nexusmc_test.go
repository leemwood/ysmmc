package nexusmc

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestSiteAPIKeyWhitespace(t *testing.T) {
	t.Setenv("NEXUSMC_SITE_API_KEY", " \t\r\n")
	if SiteConfigured() {
		t.Fatal("whitespace is not a configured key")
	}
	_, status, err := SiteGet("/test-key", nil, nil, 0)
	if status != http.StatusServiceUnavailable || !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected missing configuration, got %d %v", status, err)
	}
}

func TestSiteAPIKeyHeaderAndRecovery(t *testing.T) {
	requests := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Header.Get("X-API-Key") != "avm_sak_regression" {
			t.Error("site key was not normalized into the X-API-Key header")
		}
		w.Header().Set("Content-Type", "application/json")
		if requests == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"code":"SITE_API_KEY_INVALID"}`))
			return
		}
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer upstream.Close()
	t.Setenv("NEXUSMC_SITE_ORIGIN", upstream.URL)
	t.Setenv("NEXUSMC_SITE_API_KEY", " \tavm_sak_regression\r\n")
	const path = "/test-key-recovery"
	cache.Delete(path + "?")
	t.Cleanup(func() { cache.Delete(path + "?") })
	for i, want := range []int{http.StatusUnauthorized, http.StatusOK, http.StatusOK} {
		_, status, err := SiteGet(path, nil, nil, time.Minute)
		if err != nil || status != want {
			t.Fatalf("request %d: expected HTTP %d, got %d %v", i+1, want, status, err)
		}
	}
	if requests != 2 {
		t.Fatalf("credential failure must not be cached; success should be cached: %d upstream requests", requests)
	}
}

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
