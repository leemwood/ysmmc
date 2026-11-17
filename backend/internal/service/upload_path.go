package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
)

var ErrInvalidUpload = errors.New("无效的上传文件路径")

// ResolveUploadPath 兼容已有 /uploads/... 和相对路径，但不允许越界或符号链接。
// 模型磁盘读取和删除均使用同一个校验入口。
func ResolveUploadPath(storedPath, category string) (string, error) {
	rel := strings.ReplaceAll(storedPath, "\\", "/")
	rel = strings.TrimPrefix(rel, "/uploads/")
	if !strings.HasPrefix(rel, category+"/") || strings.Contains(rel, ":") {
		return "", ErrInvalidUpload
	}
	for _, part := range strings.Split(rel, "/") {
		if part == "" || part == "." || part == ".." {
			return "", ErrInvalidUpload
		}
	}
	root, err := filepath.Abs(config.AppConfig.UploadPath)
	if err != nil {
		return "", err
	}
	fullPath := root
	for _, part := range strings.Split(rel, "/") {
		fullPath = filepath.Join(fullPath, part)
		info, err := os.Lstat(fullPath)
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 || info.Mode()&os.ModeIrregular != 0 {
			return "", ErrInvalidUpload
		}
	}
	info, err := os.Stat(fullPath)
	if err != nil || !info.Mode().IsRegular() {
		return "", ErrInvalidUpload
	}
	return fullPath, nil
}

func removeUpload(storedPath, category string) error {
	if storedPath == "" {
		return nil
	}
	fullPath, err := ResolveUploadPath(storedPath, category)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.Remove(fullPath)
}

func validateModelUpload(userID uuid.UUID, storedPath string) (string, int64, error) {
	if _, err := ResolveUploadPath(storedPath, "models"); err != nil {
		return "", 0, ErrInvalidUpload
	}
	canonical := "/uploads/" + strings.TrimPrefix(strings.ReplaceAll(storedPath, "\\", "/"), "/uploads/")
	var upload model.ModelUpload
	if err := database.DB.Where("path = ? AND user_id = ?", canonical, userID).First(&upload).Error; err != nil {
		return "", 0, fmt.Errorf("%w：模型文件不存在或不属于当前用户，请重新上传", ErrInvalidUpload)
	}
	return canonical, upload.Size, nil
}
