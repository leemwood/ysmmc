package service

import (
	"errors"
	"io"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
	"gorm.io/gorm"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) SaveFile(name string, mimeType string, data []byte, category string, userID *uuid.UUID) (*model.File, error) {
	if !model.IsValidImageMimeType(mimeType) {
		return nil, errors.New("invalid file type, only images are allowed")
	}

	maxSize := int64(10 * 1024 * 1024)
	if int64(len(data)) > maxSize {
		return nil, errors.New("file too large, maximum size is 10MB")
	}

	file := &model.File{
		Name:     name,
		MimeType: mimeType,
		Size:     int64(len(data)),
		Data:     data,
		Category: category,
		UserID:   userID,
	}

	if err := database.DB.Create(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (s *FileService) SaveFileFromReader(name string, mimeType string, reader io.Reader, category string, userID *uuid.UUID) (*model.File, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return s.SaveFile(name, mimeType, data, category, userID)
}

func (s *FileService) GetFile(id uuid.UUID) (*model.File, error) {
	var file model.File
	if err := database.DB.First(&file, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return &file, nil
}

func (s *FileService) GetFileData(id uuid.UUID) ([]byte, string, error) {
	file, err := s.GetFile(id)
	if err != nil {
		return nil, "", err
	}
	return file.Data, file.MimeType, nil
}

func (s *FileService) DeleteFile(id uuid.UUID) error {
	result := database.DB.Delete(&model.File{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("file not found")
	}
	return nil
}

func (s *FileService) GetFilesByUser(userID uuid.UUID, category string) ([]model.File, error) {
	var files []model.File
	query := database.DB.Where("user_id = ?", userID)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (s *FileService) GetFilesByIDs(ids []uuid.UUID) ([]model.File, error) {
	var files []model.File
	if err := database.DB.Where("id IN ?", ids).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
