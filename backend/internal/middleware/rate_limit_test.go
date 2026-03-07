package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRateLimit_Allowed(t *testing.T) {
	r := gin.New()
	r.Use(RateLimit(RateLimitConfig{
		Requests: 3,
		Window:   time.Minute,
		KeyFunc:  IPKeyFunc,
	}))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("request %d should be allowed, got status %d", i+1, w.Code)
		}
	}
}

func TestRateLimit_Exceeded(t *testing.T) {
	r := gin.New()
	r.Use(RateLimit(RateLimitConfig{
		Requests: 2,
		Window:   time.Minute,
		KeyFunc:  IPKeyFunc,
	}))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:1234"
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:1234"
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("third request should be rate limited, got status %d", w.Code)
	}
}

func TestRateLimit_DifferentIPs(t *testing.T) {
	r := gin.New()
	r.Use(RateLimit(RateLimitConfig{
		Requests: 1,
		Window:   time.Minute,
		KeyFunc:  IPKeyFunc,
	}))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	ips := []string{"192.168.1.10:1234", "192.168.1.11:1234", "192.168.1.12:1234"}
	for _, ip := range ips {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("request from IP %s should be allowed, got status %d", ip, w.Code)
		}
	}
}

func TestIPKeyFunc(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "192.168.1.100:54321"
	
	key := IPKeyFunc(c)
	
	if key == "" {
		t.Error("IP key should not be empty")
	}
	if len(key) < 7 {
		t.Errorf("IP key seems too short: %q", key)
	}
}

func TestLoginRateLimit(t *testing.T) {
	r := gin.New()
	r.Use(LoginRateLimit())
	r.POST("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", nil)
		req.RemoteAddr = "192.168.1.50:1234"
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("login request %d should be allowed, got status %d", i+1, w.Code)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.RemoteAddr = "192.168.1.50:1234"
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("sixth login request should be rate limited, got status %d", w.Code)
	}
}

func TestRegisterRateLimit(t *testing.T) {
	r := gin.New()
	r.Use(RegisterRateLimit())
	r.POST("/register", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", nil)
		req.RemoteAddr = "192.168.1.60:1234"
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("register request %d should be allowed, got status %d", i+1, w.Code)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", nil)
	req.RemoteAddr = "192.168.1.60:1234"
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("fourth register request should be rate limited, got status %d", w.Code)
	}
}

func TestForgotPasswordRateLimit(t *testing.T) {
	r := gin.New()
	r.Use(ForgotPasswordRateLimit())
	r.POST("/forgot-password", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/forgot-password", nil)
		req.RemoteAddr = "192.168.1.70:1234"
		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("forgot-password request %d should be allowed, got status %d", i+1, w.Code)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forgot-password", nil)
	req.RemoteAddr = "192.168.1.70:1234"
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("fourth forgot-password request should be rate limited, got status %d", w.Code)
	}
}
