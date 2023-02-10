package user

import (
	"context"
	"go.uber.org/zap"
	"hades_backend/app/logging"
	"hades_backend/app/model/user"
	repository "hades_backend/app/repository/user"
)

type Service struct {
	repository  repository.Repository
	authService *AuthService
}

func NewService(r repository.Repository) *Service {
	return &Service{
		repository:  r,
		authService: NewAuth(r),
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*user.Login, error) {
	return s.authService.Login(ctx, email, password)
}

func (s *Service) CreateUser(ctx context.Context, user *user.User) (uint, error) { //TODO: ROLES as enum
	logger := logging.FromContext(ctx)
	logger.Info("creating user", zap.String("email", user.Email))

	user.FirstLogin = true
	user.Password = s.authService.EncodePassword(user.Password)

	id, err := s.repository.Create(ctx, user)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) UpdateUser(ctx context.Context, userId uint, user *user.User) error { //todo
	logger := logging.FromContext(ctx)
	logger.Info("updating user", zap.String("email", user.Email), zap.Uint("id", userId))
	user.ID = userId
	return s.repository.Update(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, id uint) (*user.User, error) { //todo

	logger := logging.FromContext(ctx)
	logger.Info("getting user", zap.Uint("id", id))

	return s.repository.GetByID(ctx, id)
}

func (s *Service) DeleteUser(ctx context.Context, id uint) error {
	logger := logging.FromContext(ctx)
	logger.Info("deleting user", zap.Uint("id", id))

	return s.repository.Delete(ctx, id)
}
