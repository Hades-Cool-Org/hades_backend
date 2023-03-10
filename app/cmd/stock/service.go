package stock

import (
	"context"
	"fmt"
	"hades_backend/app/logging"
	"hades_backend/app/model/stock"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetProduct(ctx context.Context, stockId uint, productId uint) (*stock.ProductData, error) {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("GettingProduct -> [ stockId: %v, productId: %v ]", stockId, productId))

	return s.repository.GetProduct(ctx, stockId, productId)
}

func (s *Service) UpdateProduct(ctx context.Context, stockId uint, productId uint, product *stock.ProductData) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("UpdatingProduct -> [ stockId: %v, productId: %v, product: %v ]", stockId, productId, product))

	return s.repository.UpdateProduct(ctx, stockId, productId, product)
}

func (s *Service) AddProductToStock(ctx context.Context, stockId uint, products []*stock.ProductData) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("AddingProductToStock -> [ stockId: %v, %v ]", stockId, len(products)))

	if len(products) == 0 {
		return nil
	}

	return s.repository.AddProductToStock(ctx, stockId, products)
}

func (s *Service) RemoveProductFromStock(ctx context.Context, stockId uint, productId uint) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("RemovingProductFromStock -> [ stockId: %v, productId: %v ]", stockId, productId))

	return s.repository.RemoveProductFromStock(ctx, stockId, productId)
}

func (s *Service) CreateStock(ctx context.Context, stock *stock.Stock) (uint, error) {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("CreatingStock -> [ %v ]", stock))

	return s.repository.Create(ctx, stock)
}

func (s *Service) GetStock(ctx context.Context, stockId uint) (*stock.Stock, error) {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("GettingStock -> [ %v ]", stockId))

	return s.repository.FindByID(ctx, stockId)
}

func (s *Service) GetStockByStoreId(ctx context.Context, storeId uint) ([]*stock.Stock, error) {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("GettingStocks -> [ %v ]", storeId))

	return s.repository.FindAllByStoreID(ctx, storeId)
}

func (s *Service) UpdateStock(ctx context.Context, stock *stock.Stock) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("UpdatingStock -> [ %v ]", stock))

	return s.repository.Update(ctx, stock)
}

func (s *Service) DeleteStock(ctx context.Context, stockId uint) error {
	l := logging.FromContext(ctx)
	l.Info(fmt.Sprintf("DeletingStock -> [ %v ]", stockId))

	return s.repository.Delete(ctx, stockId)
}

func (s *Service) GetAllStocks(ctx context.Context) ([]*stock.Stock, error) {
	l := logging.FromContext(ctx)
	l.Info("GettingAllStocks")

	return s.repository.FindAll(ctx)
}
