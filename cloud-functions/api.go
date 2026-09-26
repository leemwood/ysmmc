// EdgeOne Cloud Functions entrypoint (Framework mode, Gin).
//
// The entry file is named api.go, so EdgeOne mounts this service under the
// /api prefix: a gin route /nexusmc/resources is reachable at
// /api/nexusmc/resources on the deployed site. See docs/edgeone-deploy.md.
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ysmmc/cloud-functions/nexusmc"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	g := r.Group("/nexusmc")
	{
		g.GET("/health", health)
		g.GET("/config", config)
		g.GET("/catalog", catalog)
		g.GET("/resources", resources)
		g.GET("/resource", resourceDetail)
		g.GET("/search", search)
		g.GET("/auth/start", authStart)
		g.GET("/auth/callback", authCallback)
	}

	// EdgeOne adapts the port automatically; 9000 matches the platform examples.
	if err := r.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}

func isHTTPS(c *gin.Context) bool {
	return c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}

func proxyJSON(c *gin.Context, body []byte, status int, err error) {
	if err != nil {
		if errors.Is(err, nexusmc.ErrNotConfigured) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "not_configured",
				"message": "服务端尚未配置 NexusMC 站点 API 凭据（NEXUSMC_SITE_API_KEY）",
			})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "upstream_error",
			"message": err.Error(),
		})
		return
	}
	c.Data(status, "application/json; charset=utf-8", body)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"site_api":     nexusmc.SiteConfigured(),
		"oauth_config": nexusmc.OAuthConfigured(),
	})
}

func config(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"configured":       nexusmc.SiteConfigured(),
		"oauth_configured": nexusmc.OAuthConfigured(),
		"auth_start":       "/api/nexusmc/auth/start",
		"defaults": gin.H{
			"platform": nexusmc.ResourcePlatform(),
			"category": nexusmc.ResourceCategory(),
		},
	})
}

func catalog(c *gin.Context) {
	body, status, err := nexusmc.Catalog(c.Request.URL.Query())
	proxyJSON(c, body, status, err)
}

func resources(c *gin.Context) {
	body, status, err := nexusmc.ListResources(c.Request.URL.Query())
	proxyJSON(c, body, status, err)
}

func resourceDetail(c *gin.Context) {
	if c.Query("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_id", "message": "缺少资源 id 参数"})
		return
	}
	body, status, err := nexusmc.ResourceDetail(c.Query("id"))
	proxyJSON(c, body, status, err)
}

func search(c *gin.Context) {
	if c.Query("q") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_query", "message": "缺少搜索词 q"})
		return
	}
	body, status, err := nexusmc.SearchResources(c.Request.URL.Query())
	proxyJSON(c, body, status, err)
}

func authStart(c *gin.Context) {
	if !nexusmc.OAuthConfigured() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "oauth_not_configured",
			"message": "服务端尚未配置 NexusMC OAuth 应用（NEXUSMC_OAUTH_CLIENT_ID / NEXUSMC_OAUTH_CLIENT_SECRET）",
		})
		return
	}
	state, err := nexusmc.RandomState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	secure := isHTTPS(c)
	// The cookie is scoped to the function routes so only they can read it.
	c.SetCookie("nexusmc_oauth_state", state, 600, "/api/nexusmc", "", secure, true)
	redirectURI := nexusmc.RedirectURI(c.Request.Host, c.Request.TLS != nil, c.GetHeader("X-Forwarded-Proto"), c.GetHeader("X-Forwarded-Host"))
	http.Redirect(c.Writer, c.Request, nexusmc.AuthorizeURL(redirectURI, state), http.StatusFound)
}

func authCallback(c *gin.Context) {
	https := isHTTPS(c)
	frontend := nexusmc.FrontendOrigin(c.Request.Host, c.Request.TLS != nil, c.GetHeader("X-Forwarded-Proto"), c.GetHeader("X-Forwarded-Host"))
	fail := func(code string) {
		http.Redirect(c.Writer, c.Request, frontend+"/nexusmc/callback?error="+url.QueryEscape(code), http.StatusFound)
	}

	if e := c.Query("error"); e != "" {
		fail(e)
		return
	}

	want, err := c.Cookie("nexusmc_oauth_state")
	if err != nil {
		want = ""
	}
	got := c.Query("state")
	// Clear the state cookie regardless of the outcome.
	c.SetCookie("nexusmc_oauth_state", "", -1, "/api/nexusmc", "", https, true)
	if want == "" || got == "" || want != got {
		fail("state_mismatch")
		return
	}

	code := c.Query("code")
	if code == "" {
		fail("missing_code")
		return
	}

	redirectURI := nexusmc.RedirectURI(c.Request.Host, c.Request.TLS != nil, c.GetHeader("X-Forwarded-Proto"), c.GetHeader("X-Forwarded-Host"))
	accessToken, err := nexusmc.ExchangeCode(code, redirectURI)
	if err != nil {
		log.Printf("[nexusmc] token exchange failed: %v", err)
		fail("token_exchange")
		return
	}

	profile, err := nexusmc.UserInfo(accessToken)
	if err != nil {
		log.Printf("[nexusmc] userinfo failed: %v", err)
		fail("userinfo")
		return
	}

	user := gin.H{
		"sub":      profile["sub"],
		"uid":      profile["uid"],
		"username": profile["username"],
		"slug":     profile["slug"],
		"avatar":   nexusmc.AbsoluteURL(profile["avatar"]),
	}
	if v, ok := profile["email"]; ok && v != nil {
		user["email"] = v
	}
	if v, ok := profile["email_verified"]; ok && v != nil {
		user["email_verified"] = v
	}
	payload, err := json.Marshal(user)
	if err != nil {
		fail("internal")
		return
	}
	http.Redirect(c.Writer, c.Request,
		frontend+"/nexusmc/callback?user="+url.QueryEscape(string(payload)), http.StatusFound)
}
