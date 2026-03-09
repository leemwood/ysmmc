package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelVersion struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ModelID       uuid.UUID  `json:"model_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_model_version"`
	VersionNumber string     `json:"version_number" gorm:"size:50;not null;uniqueIndex:idx_model_version"`
	Description   *string    `json:"description" gorm:"type:text"`
	FilePath      string     `json:"file_path" gorm:"size:500;not null"`
	FileSize      int64      `json:"file_size" gorm:"default:0"`
	ImageID       *uuid.UUID `json:"image_id" gorm:"type:uuid"`
	ImageURL      *string    `json:"image_url" gorm:"type:text"`
	Changelog     *string    `json:"changelog" gorm:"type:text"`
	IsCurrent     bool       `json:"is_current" gorm:"default:false;index"`
	Downloads     int        `json:"downloads" gorm:"default:0"`
	CreatedAt     time.Time  `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time  `json:"updated_at"`

	Model   *Model `json:"model,omitempty" gorm:"foreignKey:ModelID"`
	Image   *File  `json:"image,omitempty" gorm:"foreignKey:ImageID"`
}

func (v *ModelVersion) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

func (ModelVersion) TableName() string {
	return "model_versions"
}
