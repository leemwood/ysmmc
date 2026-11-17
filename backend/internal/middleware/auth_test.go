package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/testutil"
	"github.com/ysmmc/backend/pkg/auth"
)

func TestAuthChecksCurrentAccount(t *testing.T) {
	for _, tc := range []struct {
		name, tokenRole, role string
		banned, exists        bool
		status                int
	}{
		{"active", "user", "user", false, true, 200},
		{"banned", "user", "user", true, true, 401},
		{"admin revoked", "admin", "user", false, true, 401},
		{"deleted", "user", "user", false, false, 401},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mock := testutil.Database(t)
			id := uuid.New()
			tokens, err := auth.GenerateToken(id, "old@example.com", tc.tokenRole)
			if err != nil {
				t.Fatal(err)
			}
			rows := sqlmock.NewRows([]string{"id", "role", "is_banned", "email"})
			if tc.exists {
				rows.AddRow(id, tc.role, tc.banned, "current@example.com")
			}
			mock.ExpectQuery(`SELECT .* FROM "users"`).WithArgs(id, 1).WillReturnRows(rows)
			r := gin.New()
			r.GET("/private", Auth(), func(c *gin.Context) {
				if GetEmail(c) != "current@example.com" {
					t.Error("used stale JWT profile")
				}
				c.Status(200)
			})
			req := httptest.NewRequest(http.MethodGet, "/private", nil)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.status {
				t.Fatalf("expected %d, got %d: %s", tc.status, w.Code, w.Body)
			}
		})
	}
}

func TestOptionalAuthOnlyAllowsAbsentCredentials(t *testing.T) {
	for _, tc := range []struct {
		header string
		status int
	}{{"", 200}, {"Bearer invalid", 401}} {
		r := gin.New()
		r.GET("/public", OptionalAuth(), func(c *gin.Context) { c.Status(200) })
		req := httptest.NewRequest("GET", "/public", nil)
		req.Header.Set("Authorization", tc.header)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.status {
			t.Fatalf("header %q returned %d", tc.header, w.Code)
		}
	}
}
