package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelImage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ModelID   uuid.UUID `json:"model_id" gorm:"type:uuid;not null;index;uniqueIndex:idx_model_file"`
	FileID    uuid.UUID `json:"file_id" gorm:"type:uuid;not null;uniqueIndex:idx_model_file"`
	SortOrder int       `json:"sort_order" gorm:"default:0;index"`
	CreatedAt time.Time `json:"created_at"`

	Model *Model `json:"model,omitempty" gorm:"foreignKey:ModelID"`
	File  *File  `json:"file,omitempty" gorm:"foreignKey:FileID"`
}

func (i *ModelImage) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (ModelImage) TableName() string {
	return "model_images"
}
