package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPrivateProxyNeverCaches(t *testing.T) {
	for _, status := range []int{200, 401, 403, 502} {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		proxyMe(c, []byte(`{"count":3}`), status, nil)
		if w.Code != status || w.Header().Get("Cache-Control") != "private, no-store" || w.Header().Get("Vary") != "Cookie" {
			t.Fatalf("unsafe response: %d %v", w.Code, w.Header())
		}
	}
}
