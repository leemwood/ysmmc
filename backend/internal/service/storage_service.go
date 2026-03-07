package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ysmmc/backend/internal/config"
)

type StorageService struct {
	uploadPath          string
	maxDiskUsage        int
	enableDatePartition bool
}

func NewStorageService() *StorageService {
	cfg := config.AppConfig
	return &StorageService{
		uploadPath:          cfg.UploadPath,
		maxDiskUsage:        cfg.MaxDiskUsage,
		enableDatePartition: cfg.EnableDatePartition,
	}
}

func (s *StorageService) Initialize() error {
	dirs := []string{
		filepath.Join(s.uploadPath, "models"),
		filepath.Join(s.uploadPath, "images"),
		filepath.Join(s.uploadPath, "temp"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	absPath, _ := filepath.Abs(s.uploadPath)
	fmt.Printf("Storage initialized at: %s\n", absPath)

	return nil
}

func (s *StorageService) CheckDiskSpace() error {
	return s.checkDiskSpace()
}

func (s *StorageService) SaveFile(category, filename string, content io.Reader) (string, error) {
	subDir := ""
	if s.enableDatePartition {
		subDir = time.Now().Format("2006-01")
	}

	dirPath := filepath.Join(s.uploadPath, category, subDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(dirPath, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, content); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	relativePath := filepath.Join(category, subDir, filename)
	return relativePath, nil
}

func (s *StorageService) DeleteFile(category, filename string) error {
	pattern := filepath.Join(s.uploadPath, category, "*", filename)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		directPath := filepath.Join(s.uploadPath, category, filename)
		if _, err := os.Stat(directPath); os.IsNotExist(err) {
			return errors.New("file not found")
		}
		return os.Remove(directPath)
	}

	for _, match := range matches {
		if err := os.Remove(match); err != nil {
			return err
		}
	}

	return nil
}

func (s *StorageService) GetFilePath(category, filename string) string {
	if s.enableDatePartition {
		pattern := filepath.Join(s.uploadPath, category, "*", filename)
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			return matches[0]
		}
	}
	return filepath.Join(s.uploadPath, category, filename)
}

func (s *StorageService) FileExists(category, filename string) bool {
	path := s.GetFilePath(category, filename)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (s *StorageService) GetUploadPath() string {
	return s.uploadPath
}

func (s *StorageService) GetURLPath(category, filename string) string {
	if s.enableDatePartition {
		subDir := time.Now().Format("2006-01")
		return fmt.Sprintf("/uploads/%s/%s/%s", category, subDir, filename)
	}
	return fmt.Sprintf("/uploads/%s/%s", category, filename)
}

func (s *StorageService) CleanTempFiles() error {
	tempDir := filepath.Join(s.uploadPath, "temp")
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > 24*time.Hour {
			filePath := filepath.Join(tempDir, entry.Name())
			os.Remove(filePath)
		}
	}

	return nil
}

func (s *StorageService) ValidateFilename(filename string) error {
	if filename == "" {
		return errors.New("filename cannot be empty")
	}
	if len(filename) > 255 {
		return errors.New("filename too long")
	}
	if strings.Contains(filename, "..") {
		return errors.New("invalid filename")
	}
	if strings.ContainsAny(filename, "/\\") {
		return errors.New("invalid filename")
	}
	return nil
}
