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
	"github.com/ysmmc/backend/pkg/utils"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	sessionRepo  *repository.SessionRepository
	emailService *email.EmailService
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:     repository.NewUserRepository(),
		sessionRepo:  repository.NewSessionRepository(),
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
	AccessToken        string      `json:"access_token"`
	RefreshToken       string      `json:"refresh_token"`
	ExpiresIn          int64       `json:"expires_in"`
	User               *model.User `json:"user"`
	MustChangePassword bool        `json:"must_change_password"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

func (s *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	req.Email = utils.SanitizeString(req.Email)
	req.Username = utils.SanitizeString(req.Username)

	if !utils.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	if !utils.ValidateUsername(req.Username) {
		return nil, errors.New("username can only contain letters, numbers, underscores, hyphens, and Chinese characters")
	}

	if !utils.ValidatePassword(req.Password) {
		return nil, errors.New("password must be at least 6 characters")
	}

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

	var userRole string
	var count int64
	s.userRepo.DB.Model(&model.User{}).Count(&count)
	if count == 0 {
		userRole = "super_admin"
	} else {
		userRole = "user"
	}

	user := &model.User{
		Email:         req.Email,
		PasswordHash:  passwordHash,
		Username:      req.Username,
		Role:          userRole,
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

	if user.IsBanned {
		return nil, errors.New("your account has been banned")
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	tokens, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(config.AppConfig.JWTRefreshExpireDays) * 24 * time.Hour)
	session := &model.Session{
		UserID:    user.ID,
		TokenHash: repository.HashToken(tokens.RefreshToken),
		ExpiresAt: expiresAt,
	}
	s.sessionRepo.Create(session)

	return &LoginResponse{
		AccessToken:        tokens.AccessToken,
		RefreshToken:       tokens.RefreshToken,
		ExpiresIn:          tokens.ExpiresIn,
		User:               user,
		MustChangePassword: user.MustChangePassword,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*auth.TokenPair, error) {
	tokenHash := repository.HashToken(refreshToken)
	session, err := s.sessionRepo.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		s.sessionRepo.DeleteByTokenHash(tokenHash)
		return nil, errors.New("refresh token has expired")
	}

	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.IsBanned {
		return nil, errors.New("your account has been banned")
	}

	s.sessionRepo.DeleteByTokenHash(tokenHash)

	tokens, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(config.AppConfig.JWTRefreshExpireDays) * 24 * time.Hour)
	newSession := &model.Session{
		UserID:    user.ID,
		TokenHash: repository.HashToken(tokens.RefreshToken),
		ExpiresAt: expiresAt,
	}
	s.sessionRepo.Create(newSession)

	return tokens, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := repository.HashToken(refreshToken)
	return s.sessionRepo.DeleteByTokenHash(tokenHash)
}

func (s *AuthService) ForgotPassword(emailAddr string) error {
	user, err := s.userRepo.FindByEmail(emailAddr)
	if err != nil {
		return nil
	}

	if user.IsBanned {
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

func (s *AuthService) ChangeEmail(userID uuid.UUID, newEmail string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if s.userRepo.ExistsByEmail(newEmail) {
		return errors.New("email already in use")
	}

	changeToken := uuid.New().String()
	expires := time.Now().Add(1 * time.Hour)
	user.NewEmail = &newEmail
	user.EmailChangeToken = &changeToken
	user.EmailChangeExpires = &expires

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	if s.emailService.IsConfigured() {
		verifyLink := fmt.Sprintf("%s/verify-email-change?token=%s", config.AppConfig.FrontendURL, changeToken)
		go s.emailService.Send(newEmail, "验证新邮箱 - YSM模型站",
			fmt.Sprintf(`<html><body><h2>验证新邮箱</h2><p>您正在修改邮箱地址，请点击以下链接验证新邮箱：</p><p><a href="%s">%s</a></p><p>此链接将在1小时后失效。</p></body></html>`, verifyLink, verifyLink))
	}

	return nil
}

func (s *AuthService) VerifyEmailChange(token string) error {
	user, err := s.userRepo.FindByEmailChangeToken(token)
	if err != nil {
		return errors.New("invalid email change token")
	}

	if user.NewEmail == nil {
		return errors.New("no pending email change")
	}

	if user.EmailChangeExpires == nil || user.EmailChangeExpires.Before(time.Now()) {
		user.NewEmail = nil
		user.EmailChangeToken = nil
		user.EmailChangeExpires = nil
		s.userRepo.Update(user)
		return errors.New("email change token has expired")
	}

	user.Email = *user.NewEmail
	user.NewEmail = nil
	user.EmailChangeToken = nil
	user.EmailChangeExpires = nil
	user.EmailVerified = true

	return s.userRepo.Update(user)
}
