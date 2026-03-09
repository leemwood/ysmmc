package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
)

const MaxImagesPerModel = 10

type ModelImageService struct {
	imageRepo   *repository.ModelImageRepository
	modelRepo   *repository.ModelRepository
	fileService *FileService
}

func NewModelImageService() *ModelImageService {
	return &ModelImageService{
		imageRepo:   repository.NewModelImageRepository(),
		modelRepo:   repository.NewModelRepository(),
		fileService: NewFileService(),
	}
}

type AddModelImageRequest struct {
	FileID uuid.UUID `json:"file_id" binding:"required"`
}

type ImageOrderItem struct {
	FileID    uuid.UUID `json:"file_id"`
	SortOrder int       `json:"sort_order"`
}

type UpdateImageOrderRequest struct {
	Images []ImageOrderItem `json:"images" binding:"required"`
}

func (s *ModelImageService) AddImage(modelID, userID uuid.UUID, req *AddModelImageRequest, isAdmin bool) (*model.ModelImage, error) {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return nil, err
	}

	if m.UserID != userID && !isAdmin {
		return nil, errors.New("unauthorized")
	}

	file, err := s.fileService.GetFile(req.FileID)
	if err != nil {
		return nil, errors.New("file not found")
	}

	if !file.IsImage() {
		return nil, errors.New("file is not an image")
	}

	count, err := s.imageRepo.CountByModelID(modelID)
	if err != nil {
		return nil, err
	}

	if count >= MaxImagesPerModel {
		return nil, errors.New("maximum number of images reached (10)")
	}

	exists, err := s.imageRepo.Exists(modelID, req.FileID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("image already added to this model")
	}

	image := &model.ModelImage{
		ModelID:   modelID,
		FileID:    req.FileID,
		SortOrder: int(count),
	}

	if err := s.imageRepo.Create(image); err != nil {
		return nil, err
	}

	return image, nil
}

func (s *ModelImageService) ListImages(modelID uuid.UUID) ([]model.ModelImage, error) {
	return s.imageRepo.FindByModelIDWithFile(modelID)
}

func (s *ModelImageService) DeleteImage(modelID, fileID, userID uuid.UUID, isAdmin bool) error {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return err
	}

	if m.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}

	return s.imageRepo.Delete(modelID, fileID)
}

func (s *ModelImageService) UpdateOrder(modelID, userID uuid.UUID, req *UpdateImageOrderRequest, isAdmin bool) error {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return err
	}

	if m.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}

	for _, item := range req.Images {
		image, err := s.imageRepo.FindByModelIDAndFileID(modelID, item.FileID)
		if err != nil {
			continue
		}
		image.SortOrder = item.SortOrder
		if err := s.imageRepo.Update(image); err != nil {
			return err
		}
	}

	return nil
}
