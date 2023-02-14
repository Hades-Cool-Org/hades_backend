package vendors

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"hades_backend/app/model/vendors"
	"hades_backend/app/repository"
)

type Repository interface {
	// Create creates a new vendor
	Create(ctx context.Context, vendor *vendors.Vendor) (uint, error)
	// Update updates an existing vendor
	Update(ctx context.Context, vendor *vendors.Vendor) error
	// Delete deletes an existing vendor
	Delete(ctx context.Context, id uint) error
	// GetByID returns a vendor by id
	GetByID(ctx context.Context, id uint) (*vendors.Vendor, error)
	// GetAll returns all vendors
	GetAll(ctx context.Context) ([]*vendors.Vendor, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	err := db.AutoMigrate(&Vendor{})

	if err != nil {
		panic("oops!")
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, vendor *vendors.Vendor) (uint, error) {
	model := ToModel(vendor)

	if err := m.db.Create(model).Error; err != nil {
		return 0, repository.ParseMysqlError("vendor", err)
	}
	return model.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, vendor *vendors.Vendor) error {
	model := ToModel(vendor)
	if err := m.db.Updates(model).Error; err != nil {
		return repository.ParseMysqlError("vendor", err)
	}
	return nil
}

func (m *MySqlRepository) Delete(ctx context.Context, id uint) error {
	if err := m.db.Delete(&Vendor{}, "id = ?", id).Error; err != nil {
		return repository.ParseMysqlError("vendor", err)
	}
	return nil
}

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*vendors.Vendor, error) {
	var model Vendor
	if err := m.db.Where("id = ?", id).First(&model).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, repository.ParseMysqlError("vendor", err)
	}
	return model.ToDTO(), nil
}

func (m *MySqlRepository) GetAll(ctx context.Context) ([]*vendors.Vendor, error) {
	var models []Vendor
	if err := m.db.Find(&models).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, repository.ParseMysqlError("vendor", err)
	}

	v := make([]*vendors.Vendor, len(models))

	for i, model := range models {
		v[i] = model.ToDTO()
	}

	return v, nil
}
