package repository

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
)

type FavoriteRepository struct{}

func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{}
}

func (r *FavoriteRepository) Create(favorite *model.Favorite) error {
	return database.DB.Create(favorite).Error
}

func (r *FavoriteRepository) Delete(userID, modelID uuid.UUID) error {
	return database.DB.Where("user_id = ? AND model_id = ?", userID, modelID).Delete(&model.Favorite{}).Error
}

func (r *FavoriteRepository) FindByUserAndModel(userID, modelID uuid.UUID) (*model.Favorite, error) {
	var favorite model.Favorite
	err := database.DB.Where("user_id = ? AND model_id = ?", userID, modelID).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *FavoriteRepository) Exists(userID, modelID uuid.UUID) bool {
	var count int64
	database.DB.Model(&model.Favorite{}).Where("user_id = ? AND model_id = ?", userID, modelID).Count(&count)
	return count > 0
}

func (r *FavoriteRepository) ListByUserID(userID uuid.UUID, page, pageSize int) ([]model.Favorite, int64, error) {
	var favorites []model.Favorite
	var total int64

	query := database.DB.Model(&model.Favorite{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Model.User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&favorites).Error
	return favorites, total, err
}

func (r *FavoriteRepository) CountByModel(modelID uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Favorite{}).Where("model_id = ?", modelID).Count(&count).Error
	return count, err
}
