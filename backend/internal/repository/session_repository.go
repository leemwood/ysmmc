package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
	"gorm.io/gorm"
)

type SessionRepository struct {
	DB *gorm.DB
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{DB: database.DB}
}

func (r *SessionRepository) Create(session *model.Session) error {
	return r.DB.Create(session).Error
}

func (r *SessionRepository) FindByTokenHash(tokenHash string) (*model.Session, error) {
	var session model.Session
	err := r.DB.Where("token_hash = ?", tokenHash).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteByTokenHash(tokenHash string) error {
	return r.DB.Where("token_hash = ?", tokenHash).Delete(&model.Session{}).Error
}

func (r *SessionRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.DB.Where("user_id = ?", userID).Delete(&model.Session{}).Error
}

func (r *SessionRepository) DeleteExpired() error {
	return r.DB.Where("expires_at < ?", time.Now()).Delete(&model.Session{}).Error
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
