package repository

import (
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/database"
	"github.com/ysmmc/backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{DB: database.DB}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&model.User{}, "id = ?", id).Error
}

func (r *UserRepository) List(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	r.DB.Model(&model.User{}).Count(&total)

	offset := (page - 1) * pageSize
	err := r.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) FindByResetToken(token string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("reset_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByVerificationToken(token string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("verification_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmailChangeToken(token string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email_change_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ListPendingProfiles(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.DB.Model(&model.User{}).Where("profile_status = ?", "pending_review")
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) ListByRole(role string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.DB.Model(&model.User{}).Where("role = ?", role)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) ExistsByEmail(email string) bool {
	var count int64
	r.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *UserRepository) ExistsByUsername(username string) bool {
	var count int64
	r.DB.Model(&model.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepository) CountByRole(role string) (int64, error) {
	var count int64
	err := r.DB.Model(&model.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

func (r *UserRepository) Transaction(fn func(*gorm.DB) error) error {
	return r.DB.Transaction(fn)
}
