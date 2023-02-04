package user

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"hades_backend/app/hades_errors"
	"hades_backend/app/logger"
	user2 "hades_backend/app/models/user"
	"hades_backend/app/repository/user"
	"net/http"
)

type Service struct {
	repository  user.Repository
	logger      *zap.Logger
	authService *AuthService
}

func NewService(r user.Repository) *Service {
	return &Service{
		repository:  r,
		logger:      logger.Logger,
		authService: NewAuth(r),
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*user2.Login, error) {
	return s.authService.Login(ctx, email, password)
}

func (s *Service) CreateUser(ctx context.Context, user *user2.User) error { //TODO: ROLES as enum
	s.logger.Info("creating user", zap.String("email", user.Email))

	dbUser, err := s.repository.GetByEmail(ctx, user.Email)

	if err != nil {
		return err
	}

	if dbUser != nil {
		return hades_errors.NewHadesError(errors.New("user already exists"), http.StatusConflict)
	}

	user.FirstLogin = true

	return s.repository.Create(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, user *user2.User) error { //todo
	s.logger.Info("updating user", zap.String("email", user.Email))

	return s.repository.Update(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, id string) (*user2.User, error) { //todo
	s.logger.Info("getting user", zap.String("id", id))

	return s.repository.GetByID(ctx, id)
}

func (s *Service) DeleteUser(ctx context.Context, id string) error { //todo
	s.logger.Info("deleting user", zap.String("id", id))

	return s.repository.Delete(ctx, id)
}
