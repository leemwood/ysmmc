// Package nexusmc implements a server-side client for the NexusMC open
// platform (站点 API + OAuth2). Credentials are read from environment
// variables and must never be exposed to the browser — see
// docs/reference/nexusmc-open-platform for the upstream API contract.
package nexusmc

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	DefaultSiteOrigin = "https://www.nexusmc.cn"
	// 留空 = 授权请求不带 scope，由 NexusMC 按应用登记的默认权限发放
	defaultScopes = ""
)

// ErrNotConfigured means the EdgeOne environment variables for NexusMC are missing.
var ErrNotConfigured = errors.New("nexusmc credentials are not configured")

var httpClient = &http.Client{Timeout: 15 * time.Second}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func SiteOrigin() string         { return strings.TrimRight(env("NEXUSMC_SITE_ORIGIN", DefaultSiteOrigin), "/") }
func APIKey() string             { return os.Getenv("NEXUSMC_SITE_API_KEY") }
func OAuthClientID() string      { return os.Getenv("NEXUSMC_OAUTH_CLIENT_ID") }
func OAuthClientSecret() string  { return os.Getenv("NEXUSMC_OAUTH_CLIENT_SECRET") }
func Scopes() string             { return env("NEXUSMC_OAUTH_SCOPES", defaultScopes) }
func ResourcePlatform() string   { return os.Getenv("NEXUSMC_RESOURCE_PLATFORM") }
func ResourceCategory() string   { return os.Getenv("NEXUSMC_RESOURCE_CATEGORY") }
func FrontendOriginEnv() string  { return strings.TrimRight(os.Getenv("NEXUSMC_FRONTEND_ORIGIN"), "/") }

func SiteConfigured() bool  { return APIKey() != "" }
func OAuthConfigured() bool { return OAuthClientID() != "" && OAuthClientSecret() != "" }

// RandomState returns a 128-bit hex string for the OAuth2 state parameter.
func RandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// --- site API client -------------------------------------------------------

// ResourceListParams are the query keys the NexusMC resource list endpoint
// accepts (docs: site-api-resources).
var ResourceListParams = []string{
	"page", "pageSize", "sort", "platform", "category", "subCategory",
	"sideSupport", "thirdLevel", "version", "tag", "collection",
	"officialTags", "customCategorySelections",
}

// SearchParams are accepted by /api/site/v1/search (docs: site-api-search).
var SearchParams = []string{"q", "type", "sort", "page", "pageSize", "platform", "category", "tag"}

type cachedResponse struct {
	body    []byte
	status  int
	expires time.Time
}

var cache sync.Map // string -> cachedResponse

func cacheGet(key string) ([]byte, int, bool) {
	v, ok := cache.Load(key)
	if !ok {
		return nil, 0, false
	}
	entry := v.(cachedResponse)
	if time.Now().After(entry.expires) {
		cache.Delete(key)
		return nil, 0, false
	}
	return entry.body, entry.status, true
}

func cachePut(key string, body []byte, status int, ttl time.Duration) {
	cache.Store(key, cachedResponse{body: body, status: status, expires: time.Now().Add(ttl)})
}

// SiteGet calls one NexusMC site API endpoint with the configured API key.
// Query keys are filtered to the allowlist; successful responses are cached
// for cacheTTL to conserve the key's rate limit (default 120 req/min).
func SiteGet(path string, query url.Values, allow []string, cacheTTL time.Duration) ([]byte, int, error) {
	if !SiteConfigured() {
		return nil, http.StatusServiceUnavailable, ErrNotConfigured
	}
	filtered := url.Values{}
	for _, k := range allow {
		if v := query.Get(k); v != "" {
			filtered.Set(k, v)
		}
	}
	key := path + "?" + filtered.Encode()
	if cacheTTL > 0 {
		if b, s, ok := cacheGet(key); ok {
			return b, s, nil
		}
	}
	req, err := http.NewRequest(http.MethodGet, SiteOrigin()+path+"?"+filtered.Encode(), nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Set("X-API-Key", APIKey())
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 && cacheTTL > 0 {
		cachePut(key, body, resp.StatusCode, cacheTTL)
	}
	return body, resp.StatusCode, nil
}

func withDefaults(query url.Values) url.Values {
	q := url.Values{}
	for k, vs := range query {
		q[k] = vs
	}
	if q.Get("platform") == "" {
		q.Set("platform", ResourcePlatform())
	}
	if q.Get("category") == "" {
		q.Set("category", ResourceCategory())
	}
	return q
}

// ListResources proxies GET /api/site/v1/resources.
func ListResources(query url.Values) ([]byte, int, error) {
	body, status, err := SiteGet("/api/site/v1/resources", withDefaults(query), ResourceListParams, time.Minute)
	if err != nil {
		return body, status, err
	}
	return decoratePageURLs(body), status, nil
}

// Catalog proxies GET /api/site/v1/catalog/resources (filter options).
func Catalog(query url.Values) ([]byte, int, error) {
	return SiteGet("/api/site/v1/catalog/resources", withDefaults(query), []string{"platform"}, 10*time.Minute)
}

// ResourceDetail proxies GET /api/site/v1/resources/{idOrSlug}.
func ResourceDetail(idOrSlug string) ([]byte, int, error) {
	if idOrSlug == "" {
		return nil, http.StatusBadRequest, errors.New("missing resource id")
	}
	body, status, err := SiteGet("/api/site/v1/resources/"+url.PathEscape(idOrSlug), url.Values{}, nil, time.Minute)
	if err != nil {
		return body, status, err
	}
	return decoratePageURLs(body), status, nil
}

// SearchResources proxies GET /api/site/v1/search limited to type=resources.
func SearchResources(query url.Values) ([]byte, int, error) {
	q := withDefaults(query)
	q.Set("type", "resources")
	if q.Get("q") == "" {
		return nil, http.StatusBadRequest, errors.New("missing search query")
	}
	body, status, err := SiteGet("/api/site/v1/search", q, SearchParams, 30*time.Second)
	if err != nil {
		return body, status, err
	}
	return decoratePageURLs(body), status, nil
}

// decoratePageURLs adds a stable `page_url` pointing at the public NexusMC
// resource page for every item, tolerating unknown upstream field names.
func decoratePageURLs(body []byte) []byte {
	var doc map[string]any
	if err := json.Unmarshal(body, &doc); err != nil {
		return body
	}
	decorate := func(m map[string]any) {
		pageURL := pageURLFor(m)
		if pageURL != "" {
			m["page_url"] = pageURL
		}
	}
	if items, ok := doc["items"].([]any); ok {
		for _, it := range items {
			if m, ok := it.(map[string]any); ok {
				decorate(m)
			}
		}
	}
	if res, ok := doc["resource"].(map[string]any); ok {
		decorate(res)
	}
	if _, ok := doc["page_url"]; !ok {
		if _, hasID := doc["id"]; hasID {
			decorate(doc)
		}
	}
	out, err := json.Marshal(doc)
	if err != nil {
		return body
	}
	return out
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

func pageURLFor(m map[string]any) string {
	id := str(m["slug"])
	if id == "" {
		id = str(m["id"])
	}
	if id == "" {
		return ""
	}
	return SiteOrigin() + "/resources/" + url.PathEscape(id)
}

// --- OAuth2 authorization code flow ---------------------------------------

func requestScheme(isTLS bool, protoHeader string) string {
	if isTLS || strings.EqualFold(protoHeader, "https") {
		return "https"
	}
	return "http"
}

// RedirectURI returns the OAuth2 callback address. NEXUSMC_OAUTH_REDIRECT_URI
// should be set in production so it matches the value registered on NexusMC
// verbatim. Without it the callback is derived from the forwarded request
// headers; EdgeOne's proxy hides the real host, so X-Forwarded-Host takes
// precedence over the internal Host the function actually sees.
func RedirectURI(host string, isTLS bool, protoHeader, forwardedHost string) string {
	if v := os.Getenv("NEXUSMC_OAUTH_REDIRECT_URI"); v != "" {
		return v
	}
	if forwardedHost != "" {
		return requestScheme(isTLS, protoHeader) + "://" + forwardedHost + "/api/nexusmc/auth/callback"
	}
	return requestScheme(isTLS, protoHeader) + "://" + host + "/api/nexusmc/auth/callback"
}

// FrontendOrigin is where the OAuth callback redirects the browser afterwards.
// Same precedence rules as RedirectURI: env var first, then the forwarded
// host, because the function's own Host is the platform-internal one.
func FrontendOrigin(host string, isTLS bool, protoHeader, forwardedHost string) string {
	if v := FrontendOriginEnv(); v != "" {
		return v
	}
	if forwardedHost != "" {
		return requestScheme(isTLS, protoHeader) + "://" + forwardedHost
	}
	return requestScheme(isTLS, protoHeader) + "://" + host
}

// AuthorizeURL builds the NexusMC /authorize address (docs: site-oauth2-provider).
// Scope is omitted when not configured: the server then grants the app's
// registered default scopes, which avoids invalid_scope mismatches.
func AuthorizeURL(redirectURI, state string) string {
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", OAuthClientID())
	q.Set("redirect_uri", redirectURI)
	if s := Scopes(); s != "" {
		q.Set("scope", s)
	}
	q.Set("state", state)
	return SiteOrigin() + "/api/oauth2/authorize?" + q.Encode()
}

// ExchangeCode trades an authorization code for an access token using the
// client_secret authentication method.
func ExchangeCode(code, redirectURI string) (accessToken string, err error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", OAuthClientID())
	form.Set("client_secret", OAuthClientSecret())
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest(http.MethodPost, SiteOrigin()+"/api/oauth2/token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	var parsed struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("token endpoint returned invalid JSON (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK || parsed.AccessToken == "" {
		if parsed.Error != "" {
			return "", fmt.Errorf("token endpoint error: %s", parsed.Error)
		}
		return "", fmt.Errorf("token endpoint returned HTTP %d", resp.StatusCode)
	}
	return parsed.AccessToken, nil
}

// UserInfo fetches the authorized user's profile from /api/oauth2/userinfo.
func UserInfo(accessToken string) (map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, SiteOrigin()+"/api/oauth2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo returned HTTP %d", resp.StatusCode)
	}
	var profile map[string]any
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// AbsoluteURL prefixes site-relative paths (e.g. avatars) with the site origin.
func AbsoluteURL(v any) string {
	s := str(v)
	if s == "" {
		return ""
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}
	return SiteOrigin() + s
}
