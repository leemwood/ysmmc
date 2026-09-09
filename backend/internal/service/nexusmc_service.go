package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
	"github.com/ysmmc/backend/pkg/auth"
)

const nexusmcScope = "user:basic user:email"

// NexusMCService 封装 NexusMC OAuth2 授权码流程与站点 API 代理。
// 文档：https://www.nexusmc.cn/api/oauth2 与 /api/site/v1。
type NexusMCService struct {
	repo  *repository.NexusMCRepository
	users *repository.UserRepository
	http  *http.Client
}

func NewNexusMCService() *NexusMCService {
	return &NexusMCService{
		repo:  repository.NewNexusMCRepository(),
		users: repository.NewUserRepository(),
		http:  &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *NexusMCService) Configured() bool {
	cfg := config.AppConfig
	return cfg.NexusMCClientID != "" && cfg.NexusMCClientSecret != ""
}

// AuthorizeURL 生成发起授权的重定向地址。
func (s *NexusMCService) AuthorizeURL(state string) string {
	cfg := config.AppConfig
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", cfg.NexusMCClientID)
	q.Set("redirect_uri", cfg.NexusMCRedirectURI)
	q.Set("scope", nexusmcScope)
	q.Set("state", state)
	return cfg.NexusMCOAuthBaseURL + "/authorize?" + q.Encode()
}

type nexusmcTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int64  `json:"expires_in"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (s *NexusMCService) tokenRequest(form url.Values) (*nexusmcTokenResponse, error) {
	cfg := config.AppConfig
	form.Set("client_id", cfg.NexusMCClientID)
	form.Set("client_secret", cfg.NexusMCClientSecret)

	req, err := http.NewRequest(http.MethodPost, cfg.NexusMCOAuthBaseURL+"/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nexusmc token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var tr nexusmcTokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("nexusmc token response invalid (status %d)", resp.StatusCode)
	}
	if tr.Error != "" {
		return nil, fmt.Errorf("nexusmc token error: %s: %s", tr.Error, tr.ErrorDescription)
	}
	if resp.StatusCode != http.StatusOK || tr.AccessToken == "" {
		return nil, fmt.Errorf("nexusmc token exchange failed (status %d)", resp.StatusCode)
	}
	return &tr, nil
}

type NexusMCUserInfo struct {
	Sub               string `json:"sub"`
	UID               int64  `json:"uid"`
	Username          string `json:"username"`
	Slug              string `json:"slug"`
	Avatar            string `json:"avatar"`
	Role              string `json:"role"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
}

func (s *NexusMCService) UserInfo(accessToken string) (*NexusMCUserInfo, error) {
	cfg := config.AppConfig
	req, err := http.NewRequest(http.MethodGet, cfg.NexusMCOAuthBaseURL+"/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nexusmc userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nexusmc userinfo failed (status %d)", resp.StatusCode)
	}
	var info NexusMCUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("nexusmc userinfo response invalid")
	}
	return &info, nil
}

// Revoke 主动撤销令牌；按 RFC 7009，对不存在/已过期令牌同样视为成功。
func (s *NexusMCService) Revoke(token, tokenTypeHint string) error {
	cfg := config.AppConfig
	form := url.Values{}
	form.Set("token", token)
	if tokenTypeHint != "" {
		form.Set("token_type_hint", tokenTypeHint)
	}
	req, err := http.NewRequest(http.MethodPost, cfg.NexusMCOAuthBaseURL+"/revoke", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.http.Do(req)
	if err != nil {
		return fmt.Errorf("nexusmc revoke request failed: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nexusmc revoke failed (status %d)", resp.StatusCode)
	}
	return nil
}

type NexusMCLoginResult struct {
	User       *model.User     `json:"user"`
	Tokens     *auth.TokenPair `json:"-"`
	BindingNew bool            `json:"binding_new"`
}

// HandleCallback 处理授权回调：登录模式下查找既有绑定并签发本地 JWT；
// 绑定模式（已登录用户）下新建或更新绑定。tokens 由调用方决定如何下发。
func (s *NexusMCService) HandleCallback(code string, bindUserID *uuid.UUID) (*NexusMCLoginResult, error) {
	cfg := config.AppConfig
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", cfg.NexusMCRedirectURI)

	tr, err := s.tokenRequest(form)
	if err != nil {
		return nil, err
	}

	info, err := s.UserInfo(tr.AccessToken)
	if err != nil {
		return nil, err
	}

	binding := &model.NexusMCBinding{
		Sub:           info.Sub,
		UID:           info.UID,
		Username:      info.Username,
		Slug:          info.Slug,
		Avatar:        info.Avatar,
		NexusRole:     info.Role,
		Email:         info.Email,
		EmailVerified: info.EmailVerified,
		Scope:         tr.Scope,
		AccessToken:   tr.AccessToken,
		RefreshToken:  tr.RefreshToken,
	}
	now := time.Now()
	expiresAt := now.Add(time.Duration(tr.ExpiresIn) * time.Second)
	binding.ExpiresAt = &expiresAt

	if bindUserID != nil {
		binding.UserID = *bindUserID
		if err := s.repo.Upsert(binding); err != nil {
			return nil, fmt.Errorf("failed to save nexusmc binding: %w", err)
		}
		user, err := s.users.FindByID(*bindUserID)
		if err != nil {
			return nil, err
		}
		return &NexusMCLoginResult{User: user, BindingNew: true}, nil
	}

	// 登录模式：必须有既有绑定
	existing, err := s.repo.FindBySub(info.Sub)
	if err != nil {
		return nil, fmt.Errorf("nexusmc account not bound, please login with password and bind in profile first")
	}
	user, err := s.users.FindByID(existing.UserID)
	if err != nil {
		return nil, err
	}
	if user.IsBanned {
		return nil, fmt.Errorf("account is banned")
	}
	tokens, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	existing.AccessToken = tr.AccessToken
	existing.RefreshToken = tr.RefreshToken
	existing.ExpiresAt = &expiresAt
	existing.Username = info.Username
	existing.Slug = info.Slug
	existing.Avatar = info.Avatar
	existing.NexusRole = info.Role
	if err := s.repo.UpdateTokens(existing); err != nil {
		return nil, err
	}
	return &NexusMCLoginResult{User: user, Tokens: tokens}, nil
}

// Unbind 解绑并撤销 NexusMC 侧令牌。
func (s *NexusMCService) Unbind(userID uuid.UUID) error {
	binding, err := s.repo.FindByUserID(userID)
	if err != nil {
		return fmt.Errorf("no nexusmc binding found")
	}
	if binding.RefreshToken != "" {
		_ = s.Revoke(binding.RefreshToken, "refresh_token")
	} else if binding.AccessToken != "" {
		_ = s.Revoke(binding.AccessToken, "access_token")
	}
	return s.repo.DeleteByUserID(userID)
}

// siteAPIGet 调用 NexusMC 站点 API。带上 oauthBearer 时发送双凭据（X-API-Key + Bearer）。
func (s *NexusMCService) siteAPIGet(path string, query url.Values, oauthBearer string) (int, json.RawMessage, error) {
	cfg := config.AppConfig
	if cfg.NexusMCSiteAPIKey == "" {
		return 0, nil, fmt.Errorf("NEXUSMC_SITE_API_KEY not configured")
	}

	u := cfg.NexusMCSiteAPIBaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("X-API-Key", cfg.NexusMCSiteAPIKey)
	req.Header.Set("Accept", "application/json")
	if oauthBearer != "" {
		req.Header.Set("Authorization", "Bearer "+oauthBearer)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("nexusmc site api request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	return resp.StatusCode, json.RawMessage(body), nil
}

// SiteMe 代理 GET /api/site/v1/me（双凭据：站点 API Key + 当前用户绑定的 OAuth Access Token）。
func (s *NexusMCService) SiteMe(userID uuid.UUID) (int, json.RawMessage, error) {
	binding, err := s.repo.FindByUserID(userID)
	if err != nil {
		return http.StatusForbidden, nil, fmt.Errorf("nexusmc account not bound")
	}
	accessToken := s.ensureAccessToken(binding)
	if accessToken == "" {
		return http.StatusUnauthorized, nil, fmt.Errorf("nexusmc access token invalid, please rebind")
	}
	return s.siteAPIGet("/me", nil, accessToken)
}

func (s *NexusMCService) ensureAccessToken(binding *model.NexusMCBinding) string {
	if binding.AccessToken != "" && binding.ExpiresAt != nil && time.Now().Before(binding.ExpiresAt.Add(-time.Minute)) {
		return binding.AccessToken
	}
	if binding.RefreshToken == "" {
		return ""
	}
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", binding.RefreshToken)
	tr, err := s.tokenRequest(form)
	if err != nil {
		return ""
	}
	binding.AccessToken = tr.AccessToken
	binding.RefreshToken = tr.RefreshToken
	binding.Scope = tr.Scope
	expiresAt := time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	binding.ExpiresAt = &expiresAt
	_ = s.repo.UpdateTokens(binding)
	return tr.AccessToken
}

// SiteResources 代理 /api/site/v1/resources 系列只读接口（仅 X-API-Key）。
func (s *NexusMCService) SiteResources(subPath string, query url.Values) (int, json.RawMessage, error) {
	return s.siteAPIGet("/resources"+subPath, query, "")
}

func (s *NexusMCService) GetBinding(userID uuid.UUID) (*model.NexusMCBinding, error) {
	return s.repo.FindByUserID(userID)
}
