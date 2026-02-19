package service

import (
	"errors"

	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
)

type ModelService struct {
	modelRepo *repository.ModelRepository
	userRepo  *repository.UserRepository
}

func NewModelService() *ModelService {
	return &ModelService{
		modelRepo: repository.NewModelRepository(),
		userRepo:  repository.NewUserRepository(),
	}
}

type CreateModelRequest struct {
	Title       string   `json:"title" binding:"required,max=255"`
	Description *string  `json:"description"`
	FilePath    string   `json:"file_path" binding:"required"`
	FileSize    int64    `json:"file_size"`
	ImageURL    *string  `json:"image_url"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"is_public"`
}

type UpdateModelRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	FilePath    *string  `json:"file_path"`
	FileSize    *int64   `json:"file_size"`
	ImageURL    *string  `json:"image_url"`
	Tags        []string `json:"tags"`
	IsPublic    *bool    `json:"is_public"`
}

func (s *ModelService) Create(userID uuid.UUID, req *CreateModelRequest) (*model.Model, error) {
	m := &model.Model{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		FilePath:    req.FilePath,
		FileSize:    req.FileSize,
		ImageURL:    req.ImageURL,
		Tags:        pq.StringArray(req.Tags),
		IsPublic:    req.IsPublic,
		Status:      "pending",
	}

	if err := s.modelRepo.Create(m); err != nil {
		return nil, err
	}

	return s.modelRepo.FindByID(m.ID)
}

func (s *ModelService) GetByID(id uuid.UUID) (*model.Model, error) {
	return s.modelRepo.FindByID(id)
}

func (s *ModelService) Update(modelID, userID uuid.UUID, req *UpdateModelRequest, isAdmin bool) (*model.Model, error) {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return nil, err
	}

	if m.UserID != userID && !isAdmin {
		return nil, errors.New("unauthorized")
	}

	if isAdmin {
		if req.Title != nil {
			m.Title = *req.Title
		}
		if req.Description != nil {
			m.Description = req.Description
		}
		if req.FilePath != nil {
			m.FilePath = *req.FilePath
		}
		if req.FileSize != nil {
			m.FileSize = *req.FileSize
		}
		if req.ImageURL != nil {
			m.ImageURL = req.ImageURL
		}
		if req.Tags != nil {
			m.Tags = pq.StringArray(req.Tags)
		}
		if req.IsPublic != nil {
			m.IsPublic = *req.IsPublic
		}
	} else {
		m.PendingChanges = &model.ModelPendingChanges{
			Title:       req.Title,
			Description: req.Description,
			FilePath:    req.FilePath,
			ImageURL:    req.ImageURL,
			Tags:        req.Tags,
			IsPublic:    req.IsPublic,
		}
		m.UpdateStatus = "pending_review"
	}

	if err := s.modelRepo.Update(m); err != nil {
		return nil, err
	}

	return s.modelRepo.FindByID(m.ID)
}

func (s *ModelService) Delete(modelID, userID uuid.UUID, isAdmin bool) error {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return err
	}

	if m.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}

	// Delete associated files
	if err := deleteModelFiles(m); err != nil {
		// Log error but continue with DB deletion
		// In a real production system, you might want to retry or handle this better
	}

	return s.modelRepo.Delete(modelID)
}

func deleteModelFiles(m *model.Model) error {
	cfg := config.AppConfig
	
	// Delete model file
	if m.FilePath != "" {
		// FilePath stored in DB is usually the full relative path from upload dir or absolute path
		// If it's relative to CWD (where binary runs), we can just remove it
		if err := os.Remove(m.FilePath); err != nil && !os.IsNotExist(err) {
			// Log error?
		}
	}

	// Delete image file if it's stored locally
	if m.ImageURL != nil && *m.ImageURL != "" {
		imageURL := *m.ImageURL
		// Check if it's a local file (e.g. starts with /uploads/)
		if strings.HasPrefix(imageURL, "/uploads/") {
			// Remove /uploads/ prefix
			relPath := strings.TrimPrefix(imageURL, "/uploads/")
			// Construct full path
			fullPath := filepath.Join(cfg.UploadPath, relPath)
			if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
				// Log error?
			}
		}
	}

	return nil
}

func (s *ModelService) ListPublic(page, pageSize int, search string) ([]model.Model, int64, error) {
	return s.modelRepo.ListPublic(page, pageSize, search)
}

func (s *ModelService) ListByUserID(userID uuid.UUID, page, pageSize int) ([]model.Model, int64, error) {
	return s.modelRepo.ListByUserID(userID, page, pageSize)
}

func (s *ModelService) IncrementDownloads(id uuid.UUID) error {
	return s.modelRepo.IncrementDownloads(id)
}

func (s *ModelService) Approve(modelID uuid.UUID) error {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return err
	}

	if m.UpdateStatus == "pending_review" && m.PendingChanges != nil {
		if m.PendingChanges.Title != nil {
			m.Title = *m.PendingChanges.Title
		}
		if m.PendingChanges.Description != nil {
			m.Description = m.PendingChanges.Description
		}
		if m.PendingChanges.FilePath != nil {
			m.FilePath = *m.PendingChanges.FilePath
		}
		if m.PendingChanges.ImageURL != nil {
			m.ImageURL = m.PendingChanges.ImageURL
		}
		if m.PendingChanges.Tags != nil {
			m.Tags = pq.StringArray(m.PendingChanges.Tags)
		}
		if m.PendingChanges.IsPublic != nil {
			m.IsPublic = *m.PendingChanges.IsPublic
		}
		m.PendingChanges = nil
		m.UpdateStatus = "idle"
	} else {
		m.Status = "approved"
	}

	return s.modelRepo.Update(m)
}

func (s *ModelService) Reject(modelID uuid.UUID, reason string) error {
	m, err := s.modelRepo.FindByID(modelID)
	if err != nil {
		return err
	}

	if m.UpdateStatus == "pending_review" {
		m.PendingChanges = nil
		m.UpdateStatus = "idle"
	} else {
		m.Status = "rejected"
	}

	m.RejectionReason = &reason
	return s.modelRepo.Update(m)
}

func (s *ModelService) ListPending(page, pageSize int) ([]model.Model, int64, error) {
	return s.modelRepo.ListPending(page, pageSize)
}

func (s *ModelService) ListPendingUpdates(page, pageSize int) ([]model.Model, int64, error) {
	return s.modelRepo.ListPendingUpdates(page, pageSize)
}

func (s *ModelService) Count() (int64, error) {
	return s.modelRepo.Count()
}

func (s *ModelService) CountByStatus(status string) (int64, error) {
	return s.modelRepo.CountByStatus(status)
}

func (s *ModelService) SumDownloads() (int64, error) {
	return s.modelRepo.SumDownloads()
}
