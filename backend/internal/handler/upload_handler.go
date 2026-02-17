package handler

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/middleware"
	"github.com/ysmmc/backend/pkg/response"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadModel(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "no file uploaded")
		return
	}
	defer file.Close()

	cfg := config.AppConfig
	if header.Size > cfg.MaxFileSize {
		response.BadRequest(c, "file too large")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".ysm": true, ".zip": true}
	if !allowedExts[ext] {
		response.BadRequest(c, "invalid file type, only .ysm and .zip are allowed")
		return
	}

	filename := uuid.New().String() + ext
	uploadPath := filepath.Join(cfg.UploadPath, "models", filename)

	if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
		response.InternalError(c, "failed to create upload directory")
		return
	}

	dst, err := os.Create(uploadPath)
	if err != nil {
		response.InternalError(c, "failed to create file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	response.Success(c, gin.H{
		"file_path": uploadPath,
		"file_name": header.Filename,
		"file_size": header.Size,
	})
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "no file uploaded")
		return
	}
	defer file.Close()

	cfg := config.AppConfig
	maxImageSize := cfg.MaxFileSize / 10
	if header.Size > maxImageSize {
		response.BadRequest(c, "image too large")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		response.BadRequest(c, "invalid image type")
		return
	}

	filename := uuid.New().String() + ext
	uploadPath := filepath.Join(cfg.UploadPath, "images", filename)

	if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
		response.InternalError(c, "failed to create upload directory")
		return
	}

	dst, err := os.Create(uploadPath)
	if err != nil {
		response.InternalError(c, "failed to create file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	response.Success(c, gin.H{
		"file_path": uploadPath,
		"file_name": header.Filename,
		"url":       "/uploads/images/" + filename,
	})
}

func (h *UploadHandler) DownloadModel(c *gin.Context) {
	modelID := c.Param("id")
	if modelID == "" {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)
	_ = userID

	cfg := config.AppConfig
	modelsDir := filepath.Join(cfg.UploadPath, "models")

	files, err := os.ReadDir(modelsDir)
	if err != nil {
		response.InternalError(c, "failed to read models directory")
		return
	}

	var targetFile string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), modelID) {
			targetFile = filepath.Join(modelsDir, f.Name())
			break
		}
	}

	if targetFile == "" {
		response.NotFound(c, "model file not found")
		return
	}

	c.FileAttachment(targetFile, filepath.Base(targetFile))
}

func (h *UploadHandler) ServeImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.BadRequest(c, "invalid filename")
		return
	}

	cfg := config.AppConfig
	imagePath := filepath.Join(cfg.UploadPath, "images", filename)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		response.NotFound(c, "image not found")
		return
	}

	c.File(imagePath)
}

func (h *UploadHandler) ServeUploads(c *gin.Context) {
	relativePath := c.Param("path")
	if relativePath == "" {
		response.BadRequest(c, "invalid path")
		return
	}

	cfg := config.AppConfig
	fullPath := filepath.Join(cfg.UploadPath, relativePath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		response.NotFound(c, "file not found")
		return
	}

	c.File(fullPath)
}

func init() {
	_ = time.Now
}
