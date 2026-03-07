package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ysmmc/backend/internal/config"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := config.AppConfig.AllowedOrigins
		
		allowedOrigin := ""
		for _, o := range allowedOrigins {
			if o == origin {
				allowedOrigin = origin
				break
			}
		}
		
		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
		if strings.HasPrefix(o, "*.") && strings.HasSuffix(origin, o[1:]) {
			return true
		}
	}
	return false
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		gin.DefaultWriter.Write([]byte(
			"[GIN] " + t.Format("2006/01/02 - 15:04:05") +
				" | " + c.ClientIP() +
				" | " + c.Request.Method +
				" | " + c.Request.URL.Path +
				" | " + latency.String() + "\n",
		))
	}
}
