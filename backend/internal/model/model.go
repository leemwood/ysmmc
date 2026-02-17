package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Title           string         `json:"title" gorm:"not null;size:255"`
	Description     *string        `json:"description" gorm:"type:text"`
	FilePath        string         `json:"file_path" gorm:"not null;size:500"`
	FileSize        int64          `json:"file_size" gorm:"default:0"`
	ImageURL        *string        `json:"image_url" gorm:"type:text"`
	Tags            []string       `json:"tags" gorm:"type:text[]"`
	IsPublic        bool           `json:"is_public" gorm:"default:true"`
	Status          string         `json:"status" gorm:"size:20;default:pending;index"`
	UpdateStatus    string         `json:"update_status" gorm:"size:20;default:idle"`
	PendingChanges  *ModelPendingChanges `json:"pending_changes" gorm:"type:jsonb"`
	Downloads       int            `json:"downloads" gorm:"default:0"`
	RejectionReason *string        `json:"rejection_reason" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt       time.Time      `json:"updated_at"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type ModelPendingChanges struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	FilePath    *string  `json:"file_path,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	IsPublic    *bool    `json:"is_public,omitempty"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (Model) TableName() string {
	return "models"
}
