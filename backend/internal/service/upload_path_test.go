package service

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/testutil"
)

func TestUploadPathsRejectOutsideFiles(t *testing.T) {
	testutil.Database(t)
	root := config.AppConfig.UploadPath
	if err := os.MkdirAll(filepath.Join(root, "models", "2026-11"), 0700); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "models", "2026-11", "file.ysm")
	if err := os.WriteFile(file, []byte("model"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/uploads/models/2026-11/file.ysm", "models/2026-11/file.ysm", `models\2026-11\file.ysm`} {
		resolved, err := ResolveUploadPath(path, "models")
		if err != nil || resolved != file {
			t.Fatalf("resolve %q: %q, %v", path, resolved, err)
		}
	}
	outside := filepath.Join(t.TempDir(), "sentinel")
	if err := os.WriteFile(outside, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{outside, "../sentinel", "/uploads/models/../../sentinel", "/uploads/images/file.ysm", "models/C:/sentinel", `models\..\..\sentinel`} {
		if err := deleteModelFiles(&model.Model{FilePath: path}); err == nil {
			t.Errorf("accepted %q", path)
		}
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatalf("outside sentinel changed: %v", err)
	}
	if err := deleteModelFiles(&model.Model{FilePath: "/uploads/models/2026-11/file.ysm"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Fatal("valid model file was not removed")
	}
}

func TestUploadPathsRejectSymlinks(t *testing.T) {
	testutil.Database(t)
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "keep.ysm"), []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(config.AppConfig.UploadPath, "models")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := ResolveUploadPath("models/keep.ysm", "models"); err == nil {
		t.Fatal("accepted symlink outside upload root")
	}
}

func TestModelUploadRequiresOwnership(t *testing.T) {
	for _, owned := range []bool{true, false} {
		t.Run(map[bool]string{true: "owner", false: "other user"}[owned], func(t *testing.T) {
			mock := testutil.Database(t)
			owner := uuid.New()
			if err := os.MkdirAll(filepath.Join(config.AppConfig.UploadPath, "models"), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(config.AppConfig.UploadPath, "models", "file.ysm"), []byte("model"), 0600); err != nil {
				t.Fatal(err)
			}
			rows := sqlmock.NewRows([]string{"path", "user_id", "size"})
			if owned {
				rows.AddRow("/uploads/models/file.ysm", owner, 5)
			}
			mock.ExpectQuery(`SELECT .* FROM "model_uploads" WHERE path = \$1 AND user_id = \$2`).
				WithArgs("/uploads/models/file.ysm", owner, 1).WillReturnRows(rows)
			path, size, err := validateModelUpload(owner, "models/file.ysm")
			if owned && (err != nil || path != "/uploads/models/file.ysm" || size != 5) {
				t.Fatalf("valid upload: %q, %d, %v", path, size, err)
			}
			if !owned && err == nil {
				t.Fatal("accepted another user's upload")
			}
		})
	}
}

func TestDeleteOwnedFile(t *testing.T) {
	for _, tc := range []struct {
		name         string
		owner, admin bool
		allowed      bool
	}{
		{"owner", true, false, true}, {"other user", false, false, false}, {"admin", false, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mock := testutil.Database(t)
			id, user, owner := uuid.New(), uuid.New(), uuid.New()
			if tc.owner {
				owner = user
			}
			mock.ExpectQuery(`SELECT .* FROM "files"`).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(id, owner))
			if tc.allowed {
				if tc.admin {
					mock.ExpectExec(`DELETE FROM "files" WHERE id = \$1$`).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
				} else {
					mock.ExpectExec(`DELETE FROM "files" WHERE id = \$1 AND user_id = \$2`).WithArgs(id, user).WillReturnResult(sqlmock.NewResult(0, 1))
				}
			}
			err := NewFileService().DeleteOwnedFile(id, user, tc.admin)
			if tc.allowed && err != nil {
				t.Fatal(err)
			}
			if !tc.allowed && !errors.Is(err, ErrFileForbidden) {
				t.Fatalf("expected forbidden, got %v", err)
			}
		})
	}
}
