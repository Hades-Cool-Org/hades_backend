package purchase_list

import (
	"context"
	"hades_backend/app/model"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetPurchaseList(ctx context.Context, id uint) (*model.PurchaseList, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) CreatePurchaseList(ctx context.Context, purchaseList *model.PurchaseList) (uint, error) {
	id, err := s.repository.Create(ctx, purchaseList)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) UpdatePurchaseList(ctx context.Context, purchaseListId uint, purchaseList *model.PurchaseList) error {
	purchaseList.ID = purchaseListId
	return s.repository.Update(ctx, purchaseList)
}

func (s *Service) DeletePurchaseList(ctx context.Context, purchaseListId uint) error {
	return s.repository.Delete(ctx, purchaseListId)
}

func (s *Service) GetPurchaseLists(ctx context.Context) ([]*model.PurchaseList, error) {
	return s.repository.GetAll(ctx)
}

func (s *Service) GetPurchaseListsByUserId(ctx context.Context, userId uint) ([]*model.PurchaseList, error) {
	return s.repository.GetByUserID(ctx, userId)
}
