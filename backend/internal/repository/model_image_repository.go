package repository

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
)

type ModelImageRepository struct{}

func NewModelImageRepository() *ModelImageRepository {
	return &ModelImageRepository{}
}

func (r *ModelImageRepository) Create(image *model.ModelImage) error {
	return database.DB.Create(image).Error
}

func (r *ModelImageRepository) FindByModelID(modelID uuid.UUID) ([]model.ModelImage, error) {
	var images []model.ModelImage
	err := database.DB.Where("model_id = ?", modelID).Order("sort_order ASC").Find(&images).Error
	return images, err
}

func (r *ModelImageRepository) FindByModelIDWithFile(modelID uuid.UUID) ([]model.ModelImage, error) {
	var images []model.ModelImage
	err := database.DB.Preload("File").Where("model_id = ?", modelID).Order("sort_order ASC").Find(&images).Error
	return images, err
}

func (r *ModelImageRepository) FindByModelIDAndFileID(modelID, fileID uuid.UUID) (*model.ModelImage, error) {
	var image model.ModelImage
	err := database.DB.Where("model_id = ? AND file_id = ?", modelID, fileID).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *ModelImageRepository) Delete(modelID, fileID uuid.UUID) error {
	return database.DB.Where("model_id = ? AND file_id = ?", modelID, fileID).Delete(&model.ModelImage{}).Error
}

func (r *ModelImageRepository) DeleteByModelID(modelID uuid.UUID) error {
	return database.DB.Where("model_id = ?", modelID).Delete(&model.ModelImage{}).Error
}

func (r *ModelImageRepository) CountByModelID(modelID uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ModelImage{}).Where("model_id = ?", modelID).Count(&count).Error
	return count, err
}

func (r *ModelImageRepository) Update(image *model.ModelImage) error {
	return database.DB.Save(image).Error
}

func (r *ModelImageRepository) Exists(modelID, fileID uuid.UUID) (bool, error) {
	var count int64
	err := database.DB.Model(&model.ModelImage{}).Where("model_id = ? AND file_id = ?", modelID, fileID).Count(&count).Error
	return count > 0, err
}
