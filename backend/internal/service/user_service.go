package service

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
	"github.com/ysmmc/backend/pkg/auth"
)

type UserService struct {
	userRepo  *repository.UserRepository
	modelRepo *repository.ModelRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:  repository.NewUserRepository(),
		modelRepo: repository.NewModelRepository(),
	}
}

type UpdateProfileRequest struct {
	Username  *string `json:"username"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (s *UserService) GetByID(id uuid.UUID) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetPublicProfile(id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	user.Email = ""
	return user, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, req *UpdateProfileRequest, isAdmin bool) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Username != nil && *req.Username != user.Username {
		if s.userRepo.ExistsByUsername(*req.Username) {
			return nil, errors.New("username already taken")
		}
	}

	if isAdmin {
		if req.Username != nil {
			user.Username = *req.Username
		}
		if req.Bio != nil {
			user.Bio = req.Bio
		}
		if req.AvatarURL != nil {
			user.AvatarURL = req.AvatarURL
		}
	} else {
		user.PendingChanges = &model.PendingChanges{
			Username:  req.Username,
			Bio:       req.Bio,
			AvatarURL: req.AvatarURL,
		}
		user.ProfileStatus = "pending_review"
	}

	return user, s.userRepo.Update(user)
}

func (s *UserService) ChangePassword(userID uuid.UUID, req *ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !auth.CheckPassword(req.OldPassword, user.PasswordHash) {
		return errors.New("incorrect old password")
	}

	passwordHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	return s.userRepo.Update(user)
}

func (s *UserService) List(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}

func (s *UserService) UpdateRole(userID uuid.UUID, role string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if role != "user" && role != "admin" && role != "super_admin" {
		return errors.New("invalid role")
	}

	user.Role = role
	return s.userRepo.Update(user)
}

func (s *UserService) SetRole(userID uuid.UUID, role string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.Role = role
	return s.userRepo.Update(user)
}

func (s *UserService) ApproveProfile(userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.PendingChanges != nil {
		if user.PendingChanges.Username != nil {
			user.Username = *user.PendingChanges.Username
		}
		if user.PendingChanges.Bio != nil {
			user.Bio = user.PendingChanges.Bio
		}
		if user.PendingChanges.AvatarURL != nil {
			user.AvatarURL = user.PendingChanges.AvatarURL
		}
		user.PendingChanges = nil
	}

	user.ProfileStatus = "approved"
	return s.userRepo.Update(user)
}

func (s *UserService) RejectProfile(userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.PendingChanges = nil
	user.ProfileStatus = "approved"
	return s.userRepo.Update(user)
}

func (s *UserService) ListPendingProfiles(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.ListPendingProfiles(page, pageSize)
}

func (s *UserService) Delete(userID uuid.UUID) error {
	// First, fetch all user models to delete their files
	// We fetch a large number to cover most cases. If a user has > 10000 models, this might miss some files,
	// but the DB records will still be deleted by cascade.
	// Ideally, we should loop until no more models are returned.
	page := 1
	pageSize := 1000
	for {
		models, _, err := s.modelRepo.ListByUserID(userID, page, pageSize)
		if err != nil {
			// If we fail to list models, we should probably still try to delete the user
			// but maybe log the error? For now, let's return error to be safe.
			return err
		}

		if len(models) == 0 {
			break
		}

		for _, m := range models {
			// Delete files associated with the model
			// deleteModelFiles is an unexported helper in model_service.go (same package)
			if err := deleteModelFiles(&m); err != nil {
				// Log error but continue
			}
		}

		if len(models) < pageSize {
			break
		}
		page++
	}

	// Delete user avatar if it's a local file
	user, err := s.userRepo.FindByID(userID)
	if err == nil && user.AvatarURL != nil && *user.AvatarURL != "" {
		cfg := config.AppConfig
		avatarURL := *user.AvatarURL
		if strings.HasPrefix(avatarURL, "/uploads/") {
			relPath := strings.TrimPrefix(avatarURL, "/uploads/")
			fullPath := filepath.Join(cfg.UploadPath, relPath)
			if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
				// Log error
			}
		}
	}

	// Delete user record (will cascade delete models, favorites, etc. in DB)
	return s.userRepo.Delete(userID)
}

func (s *UserService) Count() (int64, error) {
	return s.userRepo.Count()
}

func (s *UserService) GetSuperAdmin() (*model.User, error) {
	users, _, err := s.userRepo.ListByRole("super_admin", 1, 1)
	if err != nil || len(users) == 0 {
		return nil, errors.New("super admin not found")
	}
	return &users[0], nil
}

func (s *UserService) Update(user *model.User) error {
	return s.userRepo.Update(user)
}
