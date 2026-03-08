package handler

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

type UploadHandler struct {
	storageService *service.StorageService
	fileService    *service.FileService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		storageService: service.NewStorageService(),
		fileService:    service.NewFileService(),
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

	if err := validateZipMagicNumber(tee); err != nil {
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
		log.Printf("UploadImage: no file uploaded: %v", err)
		response.BadRequest(c, "no file uploaded")
		return
	}
	defer file.Close()

	ext := strings.ToLower(getExtension(header.Filename))
	mimeType := model.GetMimeTypeFromExtension(ext)
	log.Printf("UploadImage: filename=%s, ext=%s, mimeType=%s", header.Filename, ext, mimeType)

	if !model.IsValidImageMimeType(mimeType) {
		log.Printf("UploadImage: invalid image type: %s", mimeType)
		response.BadRequest(c, "invalid image type")
		return
	}

	maxImageSize := int64(10 * 1024 * 1024)
	if header.Size > maxImageSize {
		log.Printf("UploadImage: image too large: %d", header.Size)
		response.BadRequest(c, "image too large, maximum size is 10MB")
		return
	}

	buf := new(bytes.Buffer)
	tee := io.TeeReader(file, buf)

	if err := validateImageMagicNumber(tee, ext); err != nil {
		log.Printf("UploadImage: invalid image content: %v", err)
		response.BadRequest(c, "invalid image content: file does not match the declared type")
		return
	}

	data, err := io.ReadAll(io.MultiReader(buf, file))
	if err != nil {
		log.Printf("UploadImage: failed to read file: %v", err)
		response.InternalError(c, "failed to read file")
		return
	}

	userIDStr, _ := c.Get("user_id")
	var userID *uuid.UUID
	if userIDStr != nil {
		uid := userIDStr.(uuid.UUID)
		userID = &uid
	}

	category := c.DefaultPostForm("category", "model_image")

	savedFile, err := h.fileService.SaveFile(header.Filename, mimeType, data, category, userID)
	if err != nil {
		log.Printf("UploadImage: failed to save file: %v", err)
		response.InternalError(c, "failed to save file")
		return
	}

	log.Printf("UploadImage: file saved successfully, id=%s", savedFile.ID)
	response.Success(c, gin.H{
		"id":        savedFile.ID,
		"file_id":   savedFile.ID,
		"file_name": savedFile.Name,
		"size":      savedFile.Size,
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

