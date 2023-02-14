package product

import (
	"context"
	"go.uber.org/zap"
	"hades_backend/app/logging"
	"hades_backend/app/model/product"
	repository "hades_backend/app/repository/product"
)

type Service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetProduct(ctx context.Context, id uint) (*product.Product, error) {
	logger := logging.FromContext(ctx)
	logger.Info("getting product", zap.Uint("id", id))

	return s.repository.GetByID(ctx, id)
}

func (s *Service) CreateProduct(ctx context.Context, product *product.Product) (uint, error) {
	logger := logging.FromContext(ctx)
	logger.Info("creating product", zap.String("name", product.Name))
	id, err := s.repository.Create(ctx, product)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) UpdateProduct(ctx context.Context, productId uint, product *product.Product) error {
	logger := logging.FromContext(ctx)
	logger.Info("updating product", zap.String("name", product.Name), zap.Uint("id", productId))

	product.ID = productId
	return s.repository.Update(ctx, product)
}

func (s *Service) DeleteProduct(ctx context.Context, productId uint) error {
	logger := logging.FromContext(ctx)
	logger.Info("deleting product", zap.Uint("id", productId))

	return s.repository.Delete(ctx, productId)
}

func (s *Service) GetProducts(ctx context.Context) ([]*product.Product, error) {
	logger := logging.FromContext(ctx)
	logger.Info("getting products")

	return s.repository.GetAll(ctx)
}
