package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

type NexusMCHandler struct {
	svc *service.NexusMCService
}

func NewNexusMCHandler() *NexusMCHandler {
	return &NexusMCHandler{svc: service.NewNexusMCService()}
}

const nexusmcStateCookie = "nexusmc_oauth_state"

// Login 发起 NexusMC OAuth2 登录：设置 state cookie 并 302 到授权页。
func (h *NexusMCHandler) Login(c *gin.Context) {
	if !h.svc.Configured() {
		response.BadRequest(c, "nexusmc oauth is not configured")
		return
	}
	state := randomState()
	c.SetCookie(nexusmcStateCookie, state, 600, "/", "", false, true)
	c.Redirect(http.StatusFound, h.svc.AuthorizeURL(state))
}

// Bind 已登录用户发起绑定，state 里带上 bind 前缀与用户 ID。
func (h *NexusMCHandler) Bind(c *gin.Context) {
	if !h.svc.Configured() {
		response.BadRequest(c, "nexusmc oauth is not configured")
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	state := "bind:" + userID.String() + ":" + randomState()
	c.SetCookie(nexusmcStateCookie, state, 600, "/", "", false, true)
	c.Redirect(http.StatusFound, h.svc.AuthorizeURL(state))
}

// Callback 授权回调：校验 state 与 error，然后换 token。
func (h *NexusMCHandler) Callback(c *gin.Context) {
	frontend := config.AppConfig.FrontendURL
	fail := func(msg string) {
		c.Redirect(http.StatusFound, frontend+"/auth/nexusmc/callback?error="+url.QueryEscape(msg))
	}

	if errCode := c.Query("error"); errCode != "" {
		fail(errCode)
		return
	}
	code := c.Query("code")
	state := c.Query("state")
	cookieState, err := c.Cookie(nexusmcStateCookie)
	if err != nil || code == "" || state == "" || cookieState != state {
		fail("invalid_state")
		return
	}
	c.SetCookie(nexusmcStateCookie, "", -1, "/", "", false, true)

	var bindUserID *uuid.UUID
	if prefix, id, ok := splitBindState(cookieState); ok {
		if prefix != "bind" {
			fail("invalid_state")
			return
		}
		uid, err := uuid.Parse(id)
		if err != nil {
			fail("invalid_state")
			return
		}
		bindUserID = &uid
	}

	result, err := h.svc.HandleCallback(code, bindUserID)
	if err != nil {
		fail(err.Error())
		return
	}

	if bindUserID != nil {
		c.Redirect(http.StatusFound, frontend+"/auth/nexusmc/callback?bind=success")
		return
	}

	tokens := result.Tokens
	c.Redirect(http.StatusFound, frontend+"/auth/nexusmc/callback"+
		"#access_token="+url.QueryEscape(tokens.AccessToken)+
		"&refresh_token="+url.QueryEscape(tokens.RefreshToken)+
		"&expires_in="+url.QueryEscape(itoa(tokens.ExpiresIn)))
}

// Unbind 解绑当前用户的 NexusMC 账号（同时在 NexusMC 侧撤销令牌）。
func (h *NexusMCHandler) Unbind(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	if err := h.svc.Unbind(userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "nexusmc account unbound", nil)
}

// Binding 查询当前用户的绑定信息（不含令牌）。
func (h *NexusMCHandler) Binding(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	binding, err := h.svc.GetBinding(userID)
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, binding)
}

// Me 代理 NexusMC GET /api/site/v1/me（X-API-Key + Bearer 双凭据）。
func (h *NexusMCHandler) Me(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	status, body, err := h.svc.SiteMe(userID)
	if err != nil {
		response.Error(c, status, status, err.Error())
		return
	}
	c.Data(status, "application/json; charset=utf-8", body)
}

// Resources 代理 /api/site/v1/resources 系列（透传筛选参数）。
func (h *NexusMCHandler) Resources(c *gin.Context) {
	h.proxyResources(c, "")
}

func (h *NexusMCHandler) ResourceCategories(c *gin.Context) {
	h.proxyResources(c, "/categories")
}

func (h *NexusMCHandler) ResourceFilters(c *gin.Context) {
	h.proxyResources(c, "/filters")
}

func (h *NexusMCHandler) ResourceDetail(c *gin.Context) {
	h.proxyResources(c, "/"+c.Param("idOrSlug"))
}

func (h *NexusMCHandler) proxyResources(c *gin.Context, subPath string) {
	query := url.Values{}
	for k, vs := range c.Request.URL.Query() {
		for _, v := range vs {
			query.Add(k, v)
		}
	}
	status, body, err := h.svc.SiteResources(subPath, query)
	if err != nil {
		response.Error(c, http.StatusBadGateway, 502, err.Error())
		return
	}
	c.Data(status, "application/json; charset=utf-8", body)
}

func randomState() string {
	b := make([]byte, 24)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func splitBindState(state string) (string, string, bool) {
	for i := 0; i < len(state); i++ {
		if state[i] == ':' {
			rest := state[i+1:]
			for j := 0; j < len(rest); j++ {
				if rest[j] == ':' {
					return state[:i], rest[:j], true
				}
			}
			return "", "", false
		}
	}
	return "", "", false
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
