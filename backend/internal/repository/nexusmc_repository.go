package repository

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
)

type NexusMCRepository struct{}

func NewNexusMCRepository() *NexusMCRepository {
	return &NexusMCRepository{}
}

func (r *NexusMCRepository) Upsert(binding *model.NexusMCBinding) error {
	var existing model.NexusMCBinding
	err := database.DB.Where("sub = ?", binding.Sub).First(&existing).Error
	if err == nil {
		binding.ID = existing.ID
		binding.CreatedAt = existing.CreatedAt
	}
	return database.DB.Save(binding).Error
}

func (r *NexusMCRepository) FindByUserID(userID uuid.UUID) (*model.NexusMCBinding, error) {
	var binding model.NexusMCBinding
	err := database.DB.Where("user_id = ?", userID).First(&binding).Error
	if err != nil {
		return nil, err
	}
	return &binding, nil
}

func (r *NexusMCRepository) FindBySub(sub string) (*model.NexusMCBinding, error) {
	var binding model.NexusMCBinding
	err := database.DB.Where("sub = ?", sub).First(&binding).Error
	if err != nil {
		return nil, err
	}
	return &binding, nil
}

func (r *NexusMCRepository) UpdateTokens(binding *model.NexusMCBinding) error {
	return database.DB.Model(&model.NexusMCBinding{}).Where("id = ?", binding.ID).
		Updates(map[string]interface{}{
			"access_token":  binding.AccessToken,
			"refresh_token": binding.RefreshToken,
			"expires_at":    binding.ExpiresAt,
			"scope":         binding.Scope,
			"username":      binding.Username,
			"slug":          binding.Slug,
			"avatar":        binding.Avatar,
			"nexus_role":    binding.NexusRole,
		}).Error
}

func (r *NexusMCRepository) DeleteByUserID(userID uuid.UUID) error {
	return database.DB.Where("user_id = ?", userID).Delete(&model.NexusMCBinding{}).Error
}
