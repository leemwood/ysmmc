package model

import (
	"time"

	"github.com/google/uuid"
)

// ModelUpload 记录磁盘模型文件的归属；客户端不能自行声明服务器路径。
type ModelUpload struct {
	Path      string    `gorm:"primaryKey;size:500"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Size      int64     `gorm:"not null"`
	CreatedAt time.Time
}
