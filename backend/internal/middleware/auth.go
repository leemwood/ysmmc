package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
	"github.com/ysmmc/backend/pkg/auth"
	"github.com/ysmmc/backend/pkg/response"
)

func Auth() gin.HandlerFunc {
	return authenticate(true)
}

// OptionalAuth 允许匿名读取公开资源；提供凭据时必须验证当前账号状态。
func OptionalAuth() gin.HandlerFunc {
	return authenticate(false)
}

func authenticate(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" && !required {
			c.Next()
			return
		}
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}
		user, err := repository.NewUserRepository().FindByID(claims.UserID)
		if err != nil || user.IsBanned || user.Role != claims.Role {
			response.Unauthorized(c, "account unavailable or permissions changed, please log in again")
			c.Abort()
			return
		}

		c.Header("Cache-Control", "private, no-store")
		c.Set("user_id", user.ID)
		c.Set("email", user.Email)
		c.Set("role", user.Role)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || !model.IsAdmin(role.(string)) {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || !model.IsSuperAdmin(role.(string)) {
			response.Forbidden(c, "super admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) uuid.UUID {
	userID, _ := c.Get("user_id")
	return userID.(uuid.UUID)
}

func GetRole(c *gin.Context) string {
	role, _ := c.Get("role")
	return role.(string)
}

func GetEmail(c *gin.Context) string {
	email, _ := c.Get("email")
	return email.(string)
}
