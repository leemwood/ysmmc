package repository

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
)

type AnnouncementRepository struct{}

func NewAnnouncementRepository() *AnnouncementRepository {
	return &AnnouncementRepository{}
}

func (r *AnnouncementRepository) Create(announcement *model.Announcement) error {
	return database.DB.Create(announcement).Error
}

func (r *AnnouncementRepository) FindByID(id uuid.UUID) (*model.Announcement, error) {
	var announcement model.Announcement
	err := database.DB.First(&announcement, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

func (r *AnnouncementRepository) Update(announcement *model.Announcement) error {
	return database.DB.Save(announcement).Error
}

func (r *AnnouncementRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&model.Announcement{}, "id = ?", id).Error
}

func (r *AnnouncementRepository) ListActive() ([]model.Announcement, error) {
	var announcements []model.Announcement
	err := database.DB.Where("is_active = ?", true).Order("created_at DESC").Find(&announcements).Error
	return announcements, err
}

func (r *AnnouncementRepository) List(page, pageSize int) ([]model.Announcement, int64, error) {
	var announcements []model.Announcement
	var total int64

	database.DB.Model(&model.Announcement{}).Count(&total)

	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&announcements).Error
	return announcements, total, err
}
