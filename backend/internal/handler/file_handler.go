package handler

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		fileService: service.NewFileService(),
	}
}

func (h *FileHandler) GetFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	data, mimeType, err := h.fileService.GetFileData(id)
	if err != nil {
		response.NotFound(c, "file not found")
		return
	}

	c.Data(200, mimeType, data)
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "no file uploaded")
		return
	}
	defer file.Close()

	ext := strings.ToLower(getFileExtension(header.Filename))
	mimeType := model.GetMimeTypeFromExtension(ext)

	if !model.IsValidImageMimeType(mimeType) {
		response.BadRequest(c, "invalid file type, only images are allowed")
		return
	}

	buf := new(bytes.Buffer)
	tee := io.TeeReader(file, buf)

	if err := validateImageMagicNumber(tee, ext); err != nil {
		response.BadRequest(c, "invalid image content: file does not match the declared type")
		return
	}

	data, err := io.ReadAll(io.MultiReader(buf, file))
	if err != nil {
		response.InternalError(c, "failed to read file")
		return
	}

	category := c.DefaultPostForm("category", "general")
	userIDStr, _ := c.Get("user_id")
	var userID *uuid.UUID
	if userIDStr != nil {
		uid := userIDStr.(uuid.UUID)
		userID = &uid
	}

	savedFile, err := h.fileService.SaveFile(header.Filename, mimeType, data, category, userID)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	response.Success(c, gin.H{
		"id":       savedFile.ID,
		"name":     savedFile.Name,
		"mime_type": savedFile.MimeType,
		"size":     savedFile.Size,
		"category": savedFile.Category,
	})
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	if err := h.fileService.DeleteFile(id); err != nil {
		response.NotFound(c, "file not found")
		return
	}

	response.Success(c, gin.H{"message": "file deleted"})
}

func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

func getExtension(filename string) string {
	return getFileExtension(filename)
}

var errInvalidImage = errors.New("invalid image content")

func validateImageMagicNumber(reader io.Reader, ext string) error {
	buf := make([]byte, 16)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	buf = buf[:n]

	magicNumbers := map[string][][]byte{
		".jpg":  {{0xFF, 0xD8, 0xFF}},
		".jpeg": {{0xFF, 0xD8, 0xFF}},
		".png":  {{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
		".gif":  {{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}, {0x47, 0x49, 0x46, 0x38, 0x39, 0x61}},
		".webp": {{0x52, 0x49, 0x46, 0x46}},
	}

	magics, ok := magicNumbers[ext]
	if !ok {
		return nil
	}

	for _, magic := range magics {
		if len(buf) >= len(magic) {
			match := true
			for i, b := range magic {
				if buf[i] != b {
					match = false
					break
				}
			}
			if match {
				return nil
			}
		}
	}

	return errInvalidImage
}

var errInvalidZip = errors.New("invalid zip content")

func validateZipMagicNumber(reader io.Reader) error {
	buf := make([]byte, 16)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	buf = buf[:n]

	zipMagics := [][]byte{
		{0x50, 0x4B, 0x03, 0x04},
		{0x50, 0x4B, 0x05, 0x06},
		{0x50, 0x4B, 0x07, 0x08},
	}

	for _, magic := range zipMagics {
		if len(buf) >= len(magic) {
			match := true
			for i, b := range magic {
				if buf[i] != b {
					match = false
					break
				}
			}
			if match {
				return nil
			}
		}
	}

	return errInvalidZip
}
