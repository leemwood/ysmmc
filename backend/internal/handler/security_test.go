package handler_test

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/router"
	"github.com/ysmmc/backend/internal/testutil"
	"github.com/ysmmc/backend/pkg/auth"
)

func engine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	r := gin.New()
	router.Setup(r)
	return r
}

type uploadPathArg struct{ path *string }

func (arg uploadPathArg) Match(value driver.Value) bool {
	path, ok := value.(string)
	if !ok || !strings.HasPrefix(path, "/uploads/models/") || filepath.Ext(path) != ".ysm" {
		return false
	}
	*arg.path = path
	return true
}

func expectUser(mock sqlmock.Sqlmock, id uuid.UUID, role string) string {
	tokens, _ := auth.GenerateToken(id, "user@example.com", role)
	mock.ExpectQuery(`SELECT .* FROM "users"`).WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role", "email", "is_banned"}).AddRow(id, role, "user@example.com", false))
	return tokens.AccessToken
}

func expectModel(mock sqlmock.Sqlmock, id, owner uuid.UUID, status string, public bool) {
	mock.ExpectQuery(`SELECT .* FROM "models"`).WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "is_public", "title", "pending_changes"}).
			AddRow(id, owner, status, public, "Model", []byte(`{"title":"secret pending edit"}`)))
	mock.ExpectQuery(`SELECT .* FROM "users"`).WithArgs(owner).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(owner, "private@example.com"))
}

func TestModelVisibilityThroughRouter(t *testing.T) {
	for _, tc := range []struct {
		name, status, role string
		public, owner      bool
		want               int
	}{
		{"anonymous private", "approved", "", false, false, 404},
		{"anonymous pending", "pending", "", true, false, 404},
		{"anonymous rejected", "rejected", "", true, false, 404},
		{"other user private", "approved", "user", false, false, 404},
		{"owner private", "pending", "user", false, true, 200},
		{"admin private", "pending", "admin", false, false, 200},
		{"anonymous public", "approved", "", true, false, 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mock := testutil.Database(t)
			id, owner, caller := uuid.New(), uuid.New(), uuid.New()
			if tc.owner {
				caller = owner
			}
			token := ""
			if tc.role != "" {
				token = expectUser(mock, caller, tc.role)
			}
			expectModel(mock, id, owner, tc.status, tc.public)
			if tc.want == 200 {
				if tc.public && tc.status == "approved" {
					mock.ExpectExec(`UPDATE "models"`).WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorites"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				if token != "" {
					mock.ExpectQuery(`SELECT count\(\*\) FROM "favorites"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				}
			}
			req := httptest.NewRequest("GET", "/api/models/"+id.String(), nil)
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			w := httptest.NewRecorder()
			engine().ServeHTTP(w, req)
			if w.Code != tc.want {
				t.Fatalf("want %d, got %d: %s", tc.want, w.Code, w.Body)
			}
			if tc.name == "anonymous public" && (strings.Contains(w.Body.String(), "secret pending edit") || strings.Contains(w.Body.String(), "private@example.com")) {
				t.Fatal("public response leaked private fields")
			}
		})
	}
}

func TestPrivateModelSecondaryRoutes(t *testing.T) {
	for _, suffix := range []string{"/versions", "/images", "/file", "/download"} {
		t.Run(suffix, func(t *testing.T) {
			mock := testutil.Database(t)
			id, owner := uuid.New(), uuid.New()
			expectModel(mock, id, owner, "approved", false)
			method := "GET"
			if suffix == "/download" {
				method = "POST"
			}
			w := httptest.NewRecorder()
			engine().ServeHTTP(w, httptest.NewRequest(method, "/api/models/"+id.String()+suffix, nil))
			if w.Code != 404 && w.Code != 403 {
				t.Fatalf("private route returned %d: %s", w.Code, w.Body)
			}
		})
	}
}

func TestFileDeletionChecksOwnerThroughRouter(t *testing.T) {
	for _, tc := range []struct {
		name, role string
		owner      bool
		want       int
	}{{"other user", "user", false, 403}, {"owner", "user", true, 200}, {"admin", "admin", false, 200}} {
		t.Run(tc.name, func(t *testing.T) {
			mock := testutil.Database(t)
			id, user, owner := uuid.New(), uuid.New(), uuid.New()
			if tc.owner {
				owner = user
			}
			token := expectUser(mock, user, tc.role)
			mock.ExpectQuery(`SELECT .* FROM "files"`).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(id, owner))
			if tc.want == 200 {
				mock.ExpectExec(`DELETE FROM "files"`).WillReturnResult(sqlmock.NewResult(0, 1))
			}
			req := httptest.NewRequest("DELETE", "/api/files/"+id.String(), nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			engine().ServeHTTP(w, req)
			if w.Code != tc.want {
				t.Fatalf("want %d, got %d: %s", tc.want, w.Code, w.Body)
			}
		})
	}
}

func TestCreatingModelRejectsServerPath(t *testing.T) {
	mock := testutil.Database(t)
	token := expectUser(mock, uuid.New(), "user")
	// 使用越界相对路径；创建端点应在执行任何模型 INSERT 前拒绝。
	req := httptest.NewRequest("POST", "/api/models", strings.NewReader(`{"title":"attack","file_path":"../outside.ysm","image_url":"https://example.com/cover.png"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("unsafe path returned %d: %s", w.Code, w.Body)
	}
}

func TestUploadRecordsOwner(t *testing.T) {
	mock := testutil.Database(t)
	user := uuid.New()
	token := expectUser(mock, user, "user")
	config.AppConfig.MaxFileSize = 1024
	var uploadedPath string
	mock.ExpectExec(`INSERT INTO "model_uploads"`).WithArgs(
		uploadPathArg{path: &uploadedPath}, user, int64(len("test model file")), sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(0, 1))
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	part, err := form.CreateFormFile("file", "example.ysm")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("test model file"))
	_ = form.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/upload/model", &body)
	req.Header.Set("Content-Type", form.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("upload: %d %s", w.Code, w.Body)
	}
	if !strings.Contains(w.Body.String(), `/uploads/models/`) {
		t.Fatal("upload response did not use portable URL path")
	}
	var response struct {
		Data struct {
			FilePath string `json:"file_path"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Data.FilePath != uploadedPath {
		t.Fatalf("recorded %q, response %q", uploadedPath, response.Data.FilePath)
	}
}

func TestPrivateVersionRoutes(t *testing.T) {
	for _, suffix := range []string{"", "/file", "/download"} {
		t.Run(suffix, func(t *testing.T) {
			mock := testutil.Database(t)
			id, version, owner := uuid.New(), uuid.New(), uuid.New()
			mock.ExpectQuery(`SELECT .* FROM "model_versions"`).WithArgs(version, 1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "model_id"}).AddRow(version, id))
			// 版本查询预加载其模型与作者。
			mock.ExpectQuery(`SELECT .* FROM "models"`).WithArgs(id).
				WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "is_public"}).AddRow(id, owner, "approved", false))
			mock.ExpectQuery(`SELECT .* FROM "users"`).WithArgs(owner).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(owner))
			if suffix == "" {
				expectModel(mock, id, owner, "approved", false)
			}
			method := "GET"
			if suffix == "/download" {
				method = "POST"
			}
			w := httptest.NewRecorder()
			engine().ServeHTTP(w, httptest.NewRequest(method, "/api/models/"+id.String()+"/versions/"+version.String()+suffix, nil))
			if w.Code != 404 && w.Code != 403 {
				t.Fatalf("private version returned %d: %s", w.Code, w.Body)
			}
		})
	}
}

func TestUserModelsApplyVisibilityThroughRouter(t *testing.T) {
	for _, own := range []bool{true, false} {
		t.Run(fmt.Sprint(own), func(t *testing.T) {
			mock := testutil.Database(t)
			owner := uuid.New()
			token := ""
			query := `FROM "models" WHERE user_id = \$1`
			if own {
				token = expectUser(mock, owner, "user")
			} else {
				query += ` AND \(status = \$2 AND is_public = \$3\)`
			}
			count := mock.ExpectQuery(`SELECT count\(\*\) ` + query)
			list := sqlmock.NewRows([]string{"id"})
			if own {
				count.WithArgs(owner)
			} else {
				count.WithArgs(owner, "approved", true)
			}
			count.WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			selectQuery := mock.ExpectQuery(`SELECT .* ` + query + `.* LIMIT`)
			if own {
				selectQuery.WithArgs(owner, 10)
			} else {
				selectQuery.WithArgs(owner, "approved", true, 10)
			}
			selectQuery.WillReturnRows(list)
			req := httptest.NewRequest("GET", "/api/users/"+owner.String()+"/models", nil)
			if own {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			w := httptest.NewRecorder()
			engine().ServeHTTP(w, req)
			if w.Code != 200 {
				t.Fatalf("user models returned %d: %s", w.Code, w.Body)
			}
		})
	}
}

func TestRawModelFileRequiresPublishedModel(t *testing.T) {
	for _, published := range []bool{false, true} {
		t.Run(fmt.Sprint(published), func(t *testing.T) {
			mock := testutil.Database(t)
			if err := os.MkdirAll(filepath.Join(config.AppConfig.UploadPath, "models"), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(config.AppConfig.UploadPath, "models", "file.ysm"), []byte("model"), 0600); err != nil {
				t.Fatal(err)
			}
			count := 0
			if published {
				count = 1
			}
			mock.ExpectQuery(`SELECT count\(\*\) FROM "models" WHERE status = \$1 AND is_public = \$2 AND file_path IN`).
				WithArgs("approved", true, "models/file.ysm", "/uploads/models/file.ysm").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
			if !published {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "model_versions" JOIN models`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			}
			w := httptest.NewRecorder()
			engine().ServeHTTP(w, httptest.NewRequest("GET", "/api/uploads/models/file.ysm", nil))
			if published && w.Code != 200 || !published && w.Code != 404 {
				t.Fatalf("raw model: %d %s", w.Code, w.Body)
			}
		})
	}
}
