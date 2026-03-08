package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email             string          `json:"email" gorm:"uniqueIndex;not null;size:255"`
	PasswordHash      string          `json:"-" gorm:"not null;size:255"`
	Username          string          `json:"username" gorm:"uniqueIndex;not null;size:50"`
	AvatarID          *uuid.UUID      `json:"avatar_id" gorm:"type:uuid"`
	Avatar            *File           `json:"avatar,omitempty" gorm:"foreignKey:AvatarID"`
	AvatarURL         *string         `json:"avatar_url" gorm:"type:text"` // 兼容旧数据
	Bio               *string         `json:"bio" gorm:"type:text"`
	Role              string          `json:"role" gorm:"size:20;default:user"`
	ProfileStatus     string          `json:"profile_status" gorm:"size:20;default:approved"`
	PendingChanges    *PendingChanges `json:"pending_changes" gorm:"type:jsonb"`
	EmailVerified     bool            `json:"email_verified" gorm:"default:false"`
	VerificationToken *string         `json:"-" gorm:"column:verification_token;size:255"`
	ResetToken        *string         `json:"-" gorm:"size:255"`
	ResetExpires      *time.Time      `json:"-" gorm:"column:reset_token_expires"`
	NewEmail          *string         `json:"-" gorm:"size:255"`
	EmailChangeToken  *string         `json:"-" gorm:"size:255"`
	EmailChangeExpires *time.Time     `json:"-" gorm:"column:email_change_token_expires"`
	MustChangePassword bool           `json:"must_change_password" gorm:"default:false"`
	IsBanned          bool            `json:"is_banned" gorm:"default:false"`
	BannedAt          *time.Time      `json:"banned_at"`
	BannedReason      *string         `json:"banned_reason" gorm:"type:text"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type PendingChanges struct {
	Username  *string    `json:"username,omitempty"`
	Bio       *string    `json:"bio,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	AvatarID  *uuid.UUID `json:"avatar_id,omitempty"`
}

func (p *PendingChanges) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

func (p PendingChanges) Value() (driver.Value, error) {
	if p.Username == nil && p.Bio == nil && p.AvatarURL == nil && p.AvatarID == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}

func IsSuperAdmin(role string) bool {
	return role == "super_admin"
}

func IsAdmin(role string) bool {
	return role == "admin" || role == "super_admin"
}
