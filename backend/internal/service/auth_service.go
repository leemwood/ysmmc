package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/repository"
	"github.com/ysmmc/backend/pkg/auth"
	"github.com/ysmmc/backend/pkg/email"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	emailService *email.EmailService
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:     repository.NewUserRepository(),
		emailService: email.NewEmailService(),
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,min=2,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	User         *model.User `json:"user"`
}

func (s *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, errors.New("email already registered")
	}

	if s.userRepo.ExistsByUsername(req.Username) {
		return nil, errors.New("username already taken")
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:         req.Email,
		PasswordHash:  passwordHash,
		Username:      req.Username,
		Role:          "user",
		ProfileStatus: "approved",
		EmailVerified: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	if s.emailService.IsConfigured() {
		verifyToken := uuid.New().String()
		user.VerificationToken = &verifyToken
		s.userRepo.Update(user)

		verifyLink := fmt.Sprintf("%s/verify-email?token=%s", config.AppConfig.FrontendURL, verifyToken)
		go s.emailService.SendWelcome(user.Email, user.Username, verifyLink)
	}

	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	tokens, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*auth.TokenPair, error) {
	userID, err := auth.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return auth.GenerateToken(user.ID, user.Email, user.Role)
}

func (s *AuthService) ForgotPassword(emailAddr string) error {
	user, err := s.userRepo.FindByEmail(emailAddr)
	if err != nil {
		return nil
	}

	resetToken := uuid.New().String()
	user.ResetToken = &resetToken

	expires := time.Now().Add(1 * time.Hour)
	user.ResetExpires = &expires

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	if s.emailService.IsConfigured() {
		resetLink := fmt.Sprintf("%s/update-password?token=%s", config.AppConfig.FrontendURL, resetToken)
		go s.emailService.SendResetPassword(user.Email, resetLink)
	}

	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(token)
	if err != nil {
		return errors.New("invalid reset token")
	}

	if user.ResetExpires == nil || user.ResetExpires.Before(time.Now()) {
		return errors.New("reset token has expired")
	}

	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	user.ResetToken = nil
	user.ResetExpires = nil

	return s.userRepo.Update(user)
}

func (s *AuthService) VerifyEmail(token string) error {
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return errors.New("invalid verification token")
	}

	user.EmailVerified = true
	user.VerificationToken = nil

	return s.userRepo.Update(user)
}
