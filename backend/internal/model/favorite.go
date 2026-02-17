package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Favorite struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_model"`
	ModelID   uuid.UUID `json:"model_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_model;index"`
	CreatedAt time.Time `json:"created_at"`

	Model *Model `json:"model,omitempty" gorm:"foreignKey:ModelID"`
}

func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

func (Favorite) TableName() string {
	return "favorites"
}
