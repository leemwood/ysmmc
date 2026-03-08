package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string     `json:"name" gorm:"not null;size:255"`
	MimeType  string     `json:"mime_type" gorm:"not null;size:100"`
	Size      int64      `json:"size" gorm:"not null"`
	Data      []byte     `json:"-" gorm:"type:bytea;not null"`
	Category  string     `json:"category" gorm:"not null;size:50;index"`
	UserID    *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	CreatedAt time.Time  `json:"created_at" gorm:"index"`
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

func (File) TableName() string {
	return "files"
}

func (f *File) IsImage() bool {
	switch f.MimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

func GetAllowedImageMimeTypes() []string {
	return []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
	}
}

func IsValidImageMimeType(mimeType string) bool {
	for _, allowed := range GetAllowedImageMimeTypes() {
		if allowed == mimeType {
			return true
		}
	}
	return false
}

func GetMimeTypeFromExtension(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
