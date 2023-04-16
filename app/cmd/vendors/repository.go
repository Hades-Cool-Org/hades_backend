package vendors

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
	"hades_backend/app/model"
	"os"
)

type Repository interface {
	// Create creates a new vendor
	Create(ctx context.Context, vendor *model.Vendor) (uint, error)
	// Update updates an existing vendor
	Update(ctx context.Context, vendor *model.Vendor) error
	// Delete deletes an existing vendor
	Delete(ctx context.Context, id uint) error
	// GetByID returns a vendor by id
	GetByID(ctx context.Context, id uint) (*model.Vendor, error)
	// GetAll returns all vendors
	GetAll(ctx context.Context) ([]*model.Vendor, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	err := db.AutoMigrate(&Vendor{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, vendor *model.Vendor) (uint, error) {
	mm := ToModel(vendor)

	if err := m.db.Create(mm).Error; err != nil {
		return 0, cmd.ParseMysqlError(ctx, "vendor", err)
	}
	return mm.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, vendor *model.Vendor) error {
	mm := ToModel(vendor)
	if err := m.db.Updates(mm).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "vendor", err)
	}
	return nil
}

func (m *MySqlRepository) Delete(ctx context.Context, id uint) error {
	if err := m.db.Delete(&Vendor{}, "id = ?", id).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "vendor", err)
	}
	return nil
}

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*model.Vendor, error) {
	var mm Vendor
	if err := m.db.Where("id = ?", id).First(&mm).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, cmd.ParseMysqlError(ctx, "vendor", err)
	}
	return mm.ToDTO(), nil
}

func (m *MySqlRepository) GetAll(ctx context.Context) ([]*model.Vendor, error) {
	var models []Vendor
	if err := m.db.Find(&models).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "vendor", err)
	}

	v := make([]*model.Vendor, len(models))

	for i, mm := range models {
		v[i] = mm.ToDTO()
	}

	return v, nil
}
