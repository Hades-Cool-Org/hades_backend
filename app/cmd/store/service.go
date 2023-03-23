package store

import (
	"context"
	"fmt"
	userRepository "hades_backend/app/cmd/user"
	"hades_backend/app/logging"
	"hades_backend/app/model"
)

type Service struct {
	repository     Repository
	userRepository userRepository.Repository
}

func NewService(repository Repository, userRepository userRepository.Repository) *Service {
	return &Service{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *Service) AddCouriers(ctx context.Context, storeId uint, users []*model.User) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("AddingCouriers -> [ storeId: %v, %v ]", storeId, len(users)))

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
	//TODO: I think we could improve that by just using Associations.append
	for _, u := range usersResult {
		storeResult.Couriers = append(storeResult.Couriers, &model.User{ID: u.ID})
	}

	return s.repository.Update(ctx, storeResult)
}

// RemoveCouriers dumb way, just get all couriers and remove the ones that are in the list TODO improve
func (s *Service) RemoveCouriers(ctx context.Context, storeId uint, users []*model.User) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("RemovingCouriers -> [ storeId: %v, %v ]", storeId, len(users)))

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

	err = s.repository.RemoveCourierFromStore(ctx, storeId, usersResult)
	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("Succesfully removed couriers -> [ storeId: %v ]", storeId))

	return nil
}

func (s *Service) CreateStore(ctx context.Context, store *model.Store) (uint, error) {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("CreatingStore -> [ %s ]", store.ToLoggableString()))

	return s.repository.Create(ctx, store)
}

func (s *Service) UpdateStore(ctx context.Context, store *model.Store) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("UpdatingStore -> [ %s ]", store.ToLoggableString()))
	return s.repository.Update(ctx, store)
}

func (s *Service) GetStore(ctx context.Context, id uint) (*model.Store, error) {
	l := logging.FromContext(ctx)
	l.Info("GettingStore")
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetAllStores(ctx context.Context) ([]*model.Store, error) {
	l := logging.FromContext(ctx)
	l.Info("GettingAllStores")
	return s.repository.GetAll(ctx)
}

func (s *Service) DeleteStore(ctx context.Context, id uint) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("DeletingStore -> [ id: %v ]", id))
	return s.repository.Delete(ctx, id)
}

func (s *Service) GetStoreByUser(ctx context.Context, userId uint) ([]*model.Store, error) {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("GettingAllStores -> [ userId: %v ]", userId))

	u, err := s.userRepository.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return s.repository.GetByUserID(ctx, u.ID)
}
