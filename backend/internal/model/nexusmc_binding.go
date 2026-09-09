package model

import (
	"time"

	"github.com/google/uuid"
)

// NexusMCBinding 关联本地用户与 NexusMC 账号（OAuth2 绑定）。
// sub 是 NexusMC 用户的内部唯一 ID，作为稳定主键；username 仅用于展示。
type NexusMCBinding struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User          User       `json:"-" gorm:"foreignKey:UserID"`
	Sub           string     `json:"sub" gorm:"uniqueIndex;not null;size:64"`
	UID           int64      `json:"uid"`
	Username      string     `json:"username" gorm:"size:100"`
	Slug          string     `json:"slug" gorm:"size:100"`
	Avatar        string     `json:"avatar" gorm:"type:text"`
	NexusRole     string     `json:"nexus_role" gorm:"size:50"`
	Email         string     `json:"-" gorm:"size:255"`
	EmailVerified bool       `json:"email_verified" gorm:"default:false"`
	Scope         string     `json:"scope" gorm:"size:255"`
	AccessToken   string     `json:"-" gorm:"type:text"`
	RefreshToken  string     `json:"-" gorm:"type:text"`
	ExpiresAt     *time.Time `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
