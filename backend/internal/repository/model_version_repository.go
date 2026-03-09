package repository

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
)

type ModelVersionRepository struct{}

func NewModelVersionRepository() *ModelVersionRepository {
	return &ModelVersionRepository{}
}

func (r *ModelVersionRepository) Create(version *model.ModelVersion) error {
	return database.DB.Create(version).Error
}

func (r *ModelVersionRepository) FindByID(id uuid.UUID) (*model.ModelVersion, error) {
	var version model.ModelVersion
	err := database.DB.First(&version, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *ModelVersionRepository) FindByIDWithModel(id uuid.UUID) (*model.ModelVersion, error) {
	var version model.ModelVersion
	err := database.DB.Preload("Model").Preload("Model.User").First(&version, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *ModelVersionRepository) FindByModelID(modelID uuid.UUID) ([]model.ModelVersion, error) {
	var versions []model.ModelVersion
	err := database.DB.Where("model_id = ?", modelID).Order("created_at DESC").Find(&versions).Error
	return versions, err
}

func (r *ModelVersionRepository) FindCurrentVersion(modelID uuid.UUID) (*model.ModelVersion, error) {
	var version model.ModelVersion
	err := database.DB.Where("model_id = ? AND is_current = ?", modelID, true).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *ModelVersionRepository) FindByModelIDAndVersion(modelID uuid.UUID, versionNumber string) (*model.ModelVersion, error) {
	var version model.ModelVersion
	err := database.DB.Where("model_id = ? AND version_number = ?", modelID, versionNumber).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *ModelVersionRepository) Update(version *model.ModelVersion) error {
	return database.DB.Save(version).Error
}

func (r *ModelVersionRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&model.ModelVersion{}, "id = ?", id).Error
}

func (r *ModelVersionRepository) SetCurrentVersion(modelID, versionID uuid.UUID) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.ModelVersion{}).Where("model_id = ?", modelID).Update("is_current", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.ModelVersion{}).Where("id = ?", versionID).Update("is_current", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Model{}).Where("id = ?", modelID).Update("current_version_id", versionID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *ModelVersionRepository) IncrementDownloads(versionID uuid.UUID) error {
	return database.DB.Model(&model.ModelVersion{}).Where("id = ?", versionID).UpdateColumn("downloads", database.DB.Raw("downloads + 1")).Error
}

func (r *ModelVersionRepository) CountByModelID(modelID uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ModelVersion{}).Where("model_id = ?", modelID).Count(&count).Error
	return count, err
}

func (r *ModelVersionRepository) DeleteByModelID(modelID uuid.UUID) error {
	return database.DB.Where("model_id = ?", modelID).Delete(&model.ModelVersion{}).Error
}

func (r *ModelVersionRepository) ClearImageReferences(modelID uuid.UUID) error {
	return database.DB.Model(&model.ModelVersion{}).Where("model_id = ?", modelID).Update("image_id", nil).Error
}

func (r *ModelVersionRepository) ExistsByVersion(modelID uuid.UUID, versionNumber string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.ModelVersion{}).Where("model_id = ? AND version_number = ?", modelID, versionNumber).Count(&count).Error
	return count > 0, err
}
