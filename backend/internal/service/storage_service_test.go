package service

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ysmmc/backend/internal/config"
)

func init() {
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("JWT_SECRET", "test_jwt_secret_key_for_testing_purposes_32")
	os.Setenv("UPLOAD_PATH", "./test_uploads")
	os.Setenv("MAX_DISK_USAGE", "90")
	os.Setenv("ENABLE_DATE_PARTITION", "false")
	config.LoadConfig()
}

func TestStorageService_Initialize(t *testing.T) {
	svc := NewStorageService()

	err := svc.Initialize()
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	dirs := []string{
		filepath.Join(svc.uploadPath, "models"),
		filepath.Join(svc.uploadPath, "images"),
		filepath.Join(svc.uploadPath, "temp"),
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory %s should exist", dir)
		}
	}

	defer os.RemoveAll(svc.uploadPath)
}

func TestStorageService_SaveFile(t *testing.T) {
	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	content := strings.NewReader("test file content")
	filename := "test_file.txt"

	savedPath, err := svc.SaveFile("models", filename, content)
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}

	if savedPath == "" {
		t.Error("savedPath should not be empty")
	}

	expectedPath := filepath.Join("models", filename)
	if savedPath != expectedPath {
		t.Errorf("savedPath = %q, want %q", savedPath, expectedPath)
	}

	fullPath := filepath.Join(svc.uploadPath, savedPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("file should exist on disk")
	}
}

func TestStorageService_SaveFile_WithDatePartition(t *testing.T) {
	os.Setenv("ENABLE_DATE_PARTITION", "true")
	config.LoadConfig()
	defer func() {
		os.Setenv("ENABLE_DATE_PARTITION", "false")
		config.LoadConfig()
	}()

	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	content := strings.NewReader("test file content")
	filename := "test_file.txt"

	savedPath, err := svc.SaveFile("models", filename, content)
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}

	if !strings.Contains(savedPath, "models") {
		t.Error("savedPath should contain 'models'")
	}
}

func TestStorageService_DeleteFile(t *testing.T) {
	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	content := strings.NewReader("test file content")
	filename := "test_file.txt"

	svc.SaveFile("models", filename, content)

	err := svc.DeleteFile("models", filename)
	if err != nil {
		t.Fatalf("DeleteFile() error = %v", err)
	}

	fullPath := filepath.Join(svc.uploadPath, "models", filename)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestStorageService_DeleteFile_NotFound(t *testing.T) {
	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	err := svc.DeleteFile("models", "nonexistent.txt")
	if err == nil {
		t.Error("DeleteFile() should return error for non-existent file")
	}
}

func TestStorageService_FileExists(t *testing.T) {
	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	filename := "test_file.txt"
	content := strings.NewReader("test content")

	if svc.FileExists("models", filename) {
		t.Error("FileExists() should return false for non-existent file")
	}

	svc.SaveFile("models", filename, content)

	if !svc.FileExists("models", filename) {
		t.Error("FileExists() should return true for existing file")
	}
}

func TestStorageService_GetFilePath(t *testing.T) {
	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	filename := "test_file.txt"
	expectedPath := filepath.Join(svc.uploadPath, "models", filename)

	path := svc.GetFilePath("models", filename)
	if path != expectedPath {
		t.Errorf("GetFilePath() = %q, want %q", path, expectedPath)
	}
}

func TestStorageService_ValidateFilename(t *testing.T) {
	svc := NewStorageService()

	tests := []struct {
		filename string
		wantErr  bool
	}{
		{"valid_file.txt", false},
		{"valid-file.txt", false},
		{"valid_file_123.txt", false},
		{"", true},
		{strings.Repeat("a", 256), true},
		{"../traversal.txt", true},
		{"path/traversal.txt", true},
		{"path\\traversal.txt", true},
	}

	for _, tt := range tests {
		err := svc.ValidateFilename(tt.filename)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateFilename(%q) error = %v, wantErr %v", tt.filename, err, tt.wantErr)
		}
	}
}

func TestStorageService_CheckDiskSpace(t *testing.T) {
	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	err := svc.CheckDiskSpace()
	if err != nil && err.Error() != "disk usage exceeds maximum allowed percentage" {
		t.Errorf("CheckDiskSpace() unexpected error: %v", err)
	}
}

func TestStorageService_CleanTempFiles(t *testing.T) {
	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	oldFile := "old_temp.txt"
	oldContent := strings.NewReader("old content")
	svc.SaveFile("temp", oldFile, oldContent)

	oldPath := filepath.Join(svc.uploadPath, "temp", oldFile)
	oldTime := time.Now().Add(-25 * time.Hour)
	os.Chtimes(oldPath, oldTime, oldTime)

	newFile := "new_temp.txt"
	newContent := strings.NewReader("new content")
	svc.SaveFile("temp", newFile, newContent)

	err := svc.CleanTempFiles()
	if err != nil {
		t.Fatalf("CleanTempFiles() error = %v", err)
	}

	if svc.FileExists("temp", oldFile) {
		t.Error("old temp file should be cleaned")
	}

	if !svc.FileExists("temp", newFile) {
		t.Error("new temp file should not be cleaned")
	}
}

func TestStorageService_SaveLargeFile(t *testing.T) {
	svc := NewStorageService()
	svc.Initialize()
	defer os.RemoveAll(svc.uploadPath)

	largeContent := bytes.NewReader(make([]byte, 1024*1024))
	filename := "large_file.bin"

	savedPath, err := svc.SaveFile("models", filename, largeContent)
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}

	if !svc.FileExists("models", filename) {
		t.Error("large file should exist")
	}

	fullPath := filepath.Join(svc.uploadPath, savedPath)
	info, _ := os.Stat(fullPath)
	if info.Size() != 1024*1024 {
		t.Errorf("file size = %d, want %d", info.Size(), 1024*1024)
	}
}
