package service

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
)

type AnnouncementService struct {
	announcementRepo *repository.AnnouncementRepository
}

func NewAnnouncementService() *AnnouncementService {
	return &AnnouncementService{
		announcementRepo: repository.NewAnnouncementRepository(),
	}
}

type CreateAnnouncementRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
}

type UpdateAnnouncementRequest struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	IsActive *bool   `json:"is_active"`
}

func (s *AnnouncementService) Create(req *CreateAnnouncementRequest) (*model.Announcement, error) {
	announcement := &model.Announcement{
		Title:    req.Title,
		Content:  req.Content,
		IsActive: true,
	}

	if err := s.announcementRepo.Create(announcement); err != nil {
		return nil, err
	}

	return announcement, nil
}

func (s *AnnouncementService) GetByID(id uuid.UUID) (*model.Announcement, error) {
	return s.announcementRepo.FindByID(id)
}

func (s *AnnouncementService) Update(id uuid.UUID, req *UpdateAnnouncementRequest) (*model.Announcement, error) {
	announcement, err := s.announcementRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		announcement.Title = *req.Title
	}
	if req.Content != nil {
		announcement.Content = *req.Content
	}
	if req.IsActive != nil {
		announcement.IsActive = *req.IsActive
	}

	if err := s.announcementRepo.Update(announcement); err != nil {
		return nil, err
	}

	return announcement, nil
}

func (s *AnnouncementService) Delete(id uuid.UUID) error {
	return s.announcementRepo.Delete(id)
}

func (s *AnnouncementService) ListActive() ([]model.Announcement, error) {
	return s.announcementRepo.ListActive()
}

func (s *AnnouncementService) List(page, pageSize int) ([]model.Announcement, int64, error) {
	return s.announcementRepo.List(page, pageSize)
}
