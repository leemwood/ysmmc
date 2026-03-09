package service

import (
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
)

type ModelVersionService struct {
	versionRepo *repository.ModelVersionRepository
	modelRepo   *repository.ModelRepository
}

func NewModelVersionService() *ModelVersionService {
	return &ModelVersionService{
		versionRepo: repository.NewModelVersionRepository(),
		modelRepo:   repository.NewModelRepository(),
	}
}

type CreateVersionRequest struct {
	VersionNumber string     `json:"version_number" binding:"required,max=50"`
	Description   *string    `json:"description"`
	FilePath      string     `json:"file_path" binding:"required"`
	FileSize      int64      `json:"file_size"`
	ImageID       *uuid.UUID `json:"image_id"`
	ImageURL      *string    `json:"image_url"`
	Changelog     *string    `json:"changelog"`
}

type UpdateVersionRequest struct {
	Description *string    `json:"description"`
	ImageID     *uuid.UUID `json:"image_id"`
	ImageURL    *string    `json:"image_url"`
	Changelog   *string    `json:"changelog"`
}

var versionNumberRegex = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func (s *ModelVersionService) CreateVersion(modelID, userID uuid.UUID, req *CreateVersionRequest) (*model.ModelVersion, error) {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return nil, err
	}

	if m.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if !versionNumberRegex.MatchString(req.VersionNumber) {
		return nil, errors.New("invalid version number format, expected format: x.y.z (e.g., 1.0.0)")
	}

	exists, err := s.versionRepo.ExistsByVersion(modelID, req.VersionNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("version number already exists")
	}

	version := &model.ModelVersion{
		ModelID:       modelID,
		VersionNumber: req.VersionNumber,
		Description:   req.Description,
		FilePath:      req.FilePath,
		FileSize:      req.FileSize,
		ImageID:       req.ImageID,
		ImageURL:      req.ImageURL,
		Changelog:     req.Changelog,
		IsCurrent:     false,
	}

	if err := s.versionRepo.Create(version); err != nil {
		return nil, err
	}

	versionCount, err := s.versionRepo.CountByModelID(modelID)
	if err != nil {
		log.Printf("Warning: failed to count versions: %v", err)
		versionCount = 1
	}

	if err := s.modelRepo.Update(&model.Model{
		ID:           modelID,
		VersionCount: int(versionCount),
	}); err != nil {
		log.Printf("Warning: failed to update model version count: %v", err)
	}

	return version, nil
}

func (s *ModelVersionService) ListVersions(modelID uuid.UUID) ([]model.ModelVersion, error) {
	return s.versionRepo.FindByModelID(modelID)
}

func (s *ModelVersionService) GetVersion(versionID uuid.UUID) (*model.ModelVersion, error) {
	return s.versionRepo.FindByIDWithModel(versionID)
}

func (s *ModelVersionService) GetVersionByNumber(modelID uuid.UUID, versionNumber string) (*model.ModelVersion, error) {
	return s.versionRepo.FindByModelIDAndVersion(modelID, versionNumber)
}

func (s *ModelVersionService) SetCurrentVersion(modelID, versionID, userID uuid.UUID, isAdmin bool) error {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return err
	}

	if m.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}

	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return err
	}

	if version.ModelID != modelID {
		return errors.New("version does not belong to this model")
	}

	return s.versionRepo.SetCurrentVersion(modelID, versionID)
}

func (s *ModelVersionService) UpdateVersion(versionID, userID uuid.UUID, req *UpdateVersionRequest, isAdmin bool) (*model.ModelVersion, error) {
	version, err := s.versionRepo.FindByIDWithModel(versionID)
	if err != nil {
		return nil, err
	}

	if version.Model == nil {
		return nil, errors.New("model not found")
	}

	if version.Model.UserID != userID && !isAdmin {
		return nil, errors.New("unauthorized")
	}

	if req.Description != nil {
		version.Description = req.Description
	}
	if req.ImageID != nil {
		version.ImageID = req.ImageID
	}
	if req.ImageURL != nil {
		version.ImageURL = req.ImageURL
	}
	if req.Changelog != nil {
		version.Changelog = req.Changelog
	}

	if err := s.versionRepo.Update(version); err != nil {
		return nil, err
	}

	return s.versionRepo.FindByID(versionID)
}

func (s *ModelVersionService) DeleteVersion(modelID, versionID, userID uuid.UUID, isAdmin bool) error {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return err
	}

	if m.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}

	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return err
	}

	if version.ModelID != modelID {
		return errors.New("version does not belong to this model")
	}

	count, err := s.versionRepo.CountByModelID(modelID)
	if err != nil {
		return err
	}

	if count <= 1 {
		return errors.New("cannot delete the last version")
	}

	if version.IsCurrent {
		return errors.New("cannot delete the current version, please set another version as current first")
	}

	if version.FilePath != "" {
		if err := os.Remove(version.FilePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: failed to delete version file: %v", err)
		}
	}

	if err := s.versionRepo.Delete(versionID); err != nil {
		return err
	}

	versionCount, err := s.versionRepo.CountByModelID(modelID)
	if err != nil {
		log.Printf("Warning: failed to count versions: %v", err)
		versionCount = 1
	}

	if err := s.modelRepo.Update(&model.Model{
		ID:           modelID,
		VersionCount: int(versionCount),
	}); err != nil {
		log.Printf("Warning: failed to update model version count: %v", err)
	}

	return nil
}

func (s *ModelVersionService) IncrementDownloads(versionID uuid.UUID) error {
	return s.versionRepo.IncrementDownloads(versionID)
}

func (s *ModelVersionService) GetCurrentVersion(modelID uuid.UUID) (*model.ModelVersion, error) {
	return s.versionRepo.FindCurrentVersion(modelID)
}
