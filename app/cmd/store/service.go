package store

import (
	"context"
	"hades_backend/app/model/store"
	"hades_backend/app/model/user"
	storeRepository "hades_backend/app/repository/store"
	userRepository "hades_backend/app/repository/user"
)

type Service struct {
	repository     storeRepository.Repository
	userRepository userRepository.Repository
}

func NewService(repository storeRepository.Repository, userRepository userRepository.Repository) *Service {
	return &Service{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *Service) AddCouriers(ctx context.Context, storeId uint, users []user.User) error {

	if len(users) == 0 {
		return nil
	}

	var ids []uint

	for _, u := range users {
		ids = append(ids, u.ID)
	}

	usersResult, err := s.userRepository.GetMultipleByIds(ctx, ids)

	if err != nil {
		return err
	}

	storeResult, err := s.GetStore(ctx, storeId)

	if err != nil {
		return err
	}

	for _, u := range usersResult {
		storeResult.Couriers = append(storeResult.Couriers, &store.User{ID: u.ID})
	}

	return s.repository.Update(ctx, storeResult)
}

func (s *Service) CreateStore(ctx context.Context, store *store.Store) (uint, error) {
	return s.repository.Create(ctx, store)
}

func (s *Service) UpdateStore(ctx context.Context, store *store.Store) error {
	return s.repository.Update(ctx, store)
}

func (s *Service) GetStore(ctx context.Context, id uint) (*store.Store, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetAllStores(ctx context.Context) ([]*store.Store, error) {
	return s.repository.GetAll(ctx)
}

func (s *Service) DeleteStore(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) GetStoreByUser(ctx context.Context, userId uint) (*store.Store, error) {
	user, err := s.userRepository.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return s.repository.GetByID(ctx, user.StoreID)
}
