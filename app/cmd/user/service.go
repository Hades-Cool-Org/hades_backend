package user

import (
	"context"
	"go.uber.org/zap"
	"hades_backend/app/logging"
	"hades_backend/app/model"
)

type Service struct {
	repository  Repository
	authService *AuthService
}

func NewService(r Repository) *Service {
	return &Service{
		repository:  r,
		authService: NewAuth(r),
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*model.Login, error) {
	return s.authService.Login(ctx, email, password)
}

func (s *Service) CreateUser(ctx context.Context, user *model.User) (uint, error) { //TODO: ROLES as enum
	logger := logging.FromContext(ctx)
	logger.Info("creating user", zap.String("email", user.Email))

	fl := true
	user.FirstLogin = &fl
	user.Password = s.authService.EncodePassword(user.Password)

	id, err := s.repository.Create(ctx, user)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) UpdateUser(ctx context.Context, userId uint, user *model.User) error { //todo
	logger := logging.FromContext(ctx)
	logger.Info("updating user", zap.String("email", user.Email), zap.Uint("id", userId))
	user.ID = userId
	return s.repository.Update(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, id uint) (*model.User, error) { //todo
	logger := logging.FromContext(ctx)
	logger.Info("getting user", zap.Uint("id", id))

	return s.repository.GetByID(ctx, id)
}

func (s *Service) DeleteUser(ctx context.Context, id uint) error {
	logger := logging.FromContext(ctx)
	logger.Info("deleting user", zap.Uint("id", id))

	return s.repository.Delete(ctx, id)
}

func (s *Service) GetUsers(ctx context.Context) ([]*model.User, error) {
	logger := logging.FromContext(ctx)
	logger.Info("getting all users")

	return s.repository.GetUsers(ctx)
}
