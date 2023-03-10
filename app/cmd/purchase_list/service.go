package purchase_list

import (
	"context"
	"hades_backend/app/model/purchase_list"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetPurchaseList(ctx context.Context, id uint) (*purchase_list.PurchaseList, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) CreatePurchaseList(ctx context.Context, purchaseList *purchase_list.PurchaseList) (uint, error) {
	id, err := s.repository.Create(ctx, purchaseList)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) UpdatePurchaseList(ctx context.Context, purchaseListId uint, purchaseList *purchase_list.PurchaseList) error {
	purchaseList.ID = purchaseListId
	return s.repository.Update(ctx, purchaseList)
}

func (s *Service) DeletePurchaseList(ctx context.Context, purchaseListId uint) error {
	return s.repository.Delete(ctx, purchaseListId)
}

func (s *Service) GetPurchaseLists(ctx context.Context) ([]*purchase_list.PurchaseList, error) {
	return s.repository.GetAll(ctx)
}

func (s *Service) GetPurchaseListsByUserId(ctx context.Context, userId uint) ([]*purchase_list.PurchaseList, error) {
	return s.repository.GetByUserID(ctx, userId)
}
