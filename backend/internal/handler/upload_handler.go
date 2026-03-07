package handler

import (
	"bytes"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
	"github.com/ysmmc/backend/pkg/utils"
)

type UploadHandler struct {
	storageService *service.StorageService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		storageService: service.NewStorageService(),
	}
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

	ext := strings.ToLower(getExtension(header.Filename))
	allowedExts := map[string]bool{".ysm": true, ".zip": true}
	if !allowedExts[ext] {
		response.BadRequest(c, "invalid file type, only .ysm and .zip are allowed")
		return
	}

	if err := h.storageService.ValidateFilename(header.Filename); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	buf := new(bytes.Buffer)
	tee := io.TeeReader(file, buf)

	if err := utils.ValidateZipMagicNumber(tee); err != nil {
		response.BadRequest(c, "invalid file content: file does not appear to be a valid archive")
		return
	}

	if err := h.storageService.CheckDiskSpace(); err != nil {
		response.BadRequest(c, "insufficient disk space")
		return
	}

	filename := uuid.New().String() + ext
	multiReader := io.MultiReader(buf, file)

	savedPath, err := h.storageService.SaveFile("models", filename, multiReader)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	response.Success(c, gin.H{
		"file_path":  "/uploads/" + savedPath,
		"file_name":  header.Filename,
		"file_size":  header.Size,
		"saved_path": savedPath,
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

	ext := strings.ToLower(getExtension(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		response.BadRequest(c, "invalid image type")
		return
	}

	if err := h.storageService.ValidateFilename(header.Filename); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	buf := new(bytes.Buffer)
	tee := io.TeeReader(file, buf)

	if err := utils.ValidateImageMagicNumber(tee, ext); err != nil {
		response.BadRequest(c, "invalid image content: file does not match the declared type")
		return
	}

	if err := h.storageService.CheckDiskSpace(); err != nil {
		response.BadRequest(c, "insufficient disk space")
		return
	}

	filename := uuid.New().String() + ext
	multiReader := io.MultiReader(buf, file)

	savedPath, err := h.storageService.SaveFile("images", filename, multiReader)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	response.Success(c, gin.H{
		"file_name":  header.Filename,
		"url":        "/uploads/" + savedPath,
		"saved_path": savedPath,
	})
}

func (h *UploadHandler) ServeImage(c *gin.Context) {
	filename := strings.TrimPrefix(c.Param("filename"), "/")
	if filename == "" {
		response.BadRequest(c, "invalid filename")
		return
	}

	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		response.BadRequest(c, "invalid filename")
		return
	}

	imagePath := h.storageService.GetFilePath("images", filename)

	if !h.storageService.FileExists("images", filename) {
		response.NotFound(c, "image not found")
		return
	}

	c.File(imagePath)
}

func (h *UploadHandler) ServeModelFile(c *gin.Context) {
	filename := strings.TrimPrefix(c.Param("filename"), "/")
	if filename == "" {
		response.BadRequest(c, "invalid filename")
		return
	}

	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		response.BadRequest(c, "invalid filename")
		return
	}

	modelPath := h.storageService.GetFilePath("models", filename)

	if !h.storageService.FileExists("models", filename) {
		response.NotFound(c, "model file not found")
		return
	}

	c.FileAttachment(modelPath, filename)
}

func (h *UploadHandler) DeleteModel(filename string) error {
	return h.storageService.DeleteFile("models", filename)
}

func (h *UploadHandler) DeleteImage(filename string) error {
	return h.storageService.DeleteFile("images", filename)
}

func getExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}
