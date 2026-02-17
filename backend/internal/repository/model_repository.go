package repository

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
)

type ModelRepository struct{}

func NewModelRepository() *ModelRepository {
	return &ModelRepository{}
}

func (r *ModelRepository) Create(m *model.Model) error {
	return database.DB.Create(m).Error
}

func (r *ModelRepository) FindByID(id uuid.UUID) (*model.Model, error) {
	var m model.Model
	err := database.DB.Preload("User").First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ModelRepository) Update(m *model.Model) error {
	return database.DB.Save(m).Error
}

func (r *ModelRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&model.Model{}, "id = ?", id).Error
}

func (r *ModelRepository) ListPublic(page, pageSize int, search string) ([]model.Model, int64, error) {
	var models []model.Model
	var total int64

	query := database.DB.Model(&model.Model{}).Where("status = ? AND is_public = ?", "approved", true)

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error
	return models, total, err
}

func (r *ModelRepository) ListByUserID(userID uuid.UUID, page, pageSize int) ([]model.Model, int64, error) {
	var models []model.Model
	var total int64

	query := database.DB.Model(&model.Model{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error
	return models, total, err
}

func (r *ModelRepository) ListPending(page, pageSize int) ([]model.Model, int64, error) {
	var models []model.Model
	var total int64

	query := database.DB.Model(&model.Model{}).Where("status = ?", "pending")
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error
	return models, total, err
}

func (r *ModelRepository) ListPendingUpdates(page, pageSize int) ([]model.Model, int64, error) {
	var models []model.Model
	var total int64

	query := database.DB.Model(&model.Model{}).Where("update_status = ?", "pending_review")
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error
	return models, total, err
}

func (r *ModelRepository) IncrementDownloads(id uuid.UUID) error {
	return database.DB.Model(&model.Model{}).Where("id = ?", id).UpdateColumn("downloads", database.DB.Raw("downloads + 1")).Error
}

func (r *ModelRepository) Count() (int64, error) {
	var count int64
	err := database.DB.Model(&model.Model{}).Count(&count).Error
	return count, err
}

func (r *ModelRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Model{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *ModelRepository) SumDownloads() (int64, error) {
	var total int64
	err := database.DB.Model(&model.Model{}).Select("COALESCE(SUM(downloads), 0)").Scan(&total).Error
	return total, err
}
