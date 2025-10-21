package service

import (
	"errors"
	"template/app"
	"template/model"
	"template/util/auth"
	"template/util/cerror"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IAuthService interface {
	Login(email, password string) (string, string, error)
	RefreshTokens(user *model.User) (string, string, error)
}

type AuthService struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewAuthService() IAuthService {
	var service IAuthService

	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = &AuthService{
			db:     db,
			logger: logger,
		}
	})

	return service
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Debugf("User not found Email = %s", email)
			return "", "", cerror.ErrInvalidCredentials
		}

		s.logger.Errorf("Failed to query user, error = %+v", err)
		return "", "", err
	}

	if !auth.VerifyPassword(user.PasswordHash, password) {
		s.logger.Debugf("Invalid password for user Email: %s, uuid: %s", user.Email, user.Uuid)
		return "", "", cerror.ErrInvalidCredentials
	}

	token, refresh, err := auth.GenerateTokens(&user)
	if err != nil {
		s.logger.Errorf("Failed to generate token error = %+v", err)
		return "", "", err
	}

	return token, refresh, nil
}

func (s *AuthService) RefreshTokens(user *model.User) (string, string, error) {
	// TODO: implement better refresh
	return auth.GenerateTokens(user)
}
