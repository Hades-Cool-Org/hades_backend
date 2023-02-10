package vendors

import (
	"context"
	"go.uber.org/zap"
	"hades_backend/app/logger"
	"hades_backend/app/model/vendors"
	repository "hades_backend/app/repository/vendors"
)

type Service struct {
	repository repository.Repository
	logger     *zap.Logger
}

func NewService(r repository.Repository) *Service {
	return &Service{
		repository: r,
		logger:     logger.Logger,
	}
}

func (s *Service) GetVendor(ctx context.Context, id uint) (*vendors.Vendor, error) {
	s.logger.Info("getting vendor", zap.Uint("id", id))

	return s.repository.GetByID(ctx, id)
}

func (s *Service) CreateVendor(ctx context.Context, vendor *vendors.Vendor) (uint, error) {
	s.logger.Info("creating vendor", zap.String("name", vendor.Name))

	id, err := s.repository.Create(ctx, vendor)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) UpdateVendor(ctx context.Context, vendorId uint, vendor *vendors.Vendor) error {
	s.logger.Info("updating vendor", zap.String("name", vendor.Name), zap.Uint("id", vendorId))
	vendor.ID = vendorId
	return s.repository.Update(ctx, vendor)
}

func (s *Service) DeleteVendor(ctx context.Context, vendorId uint) error {
	s.logger.Info("deleting vendor", zap.Uint("id", vendorId))

	return s.repository.Delete(ctx, vendorId)
}

func (s *Service) GetVendors(ctx context.Context) ([]*vendors.Vendor, error) {
	s.logger.Info("getting vendors")

	return s.repository.GetAll(ctx)
}