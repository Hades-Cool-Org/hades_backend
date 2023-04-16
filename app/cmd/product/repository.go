package product

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
	"hades_backend/app/model"
	"os"
)

type Repository interface {
	// Create creates a new product
	Create(ctx context.Context, product *model.Product) (uint, error)
	// Update updates an existing product
	Update(ctx context.Context, product *model.Product) error
	// Delete deletes an existing product
	Delete(ctx context.Context, id uint) error
	// GetByID returns a product by id
	GetByID(ctx context.Context, id uint) (*model.Product, error)
	// GetAll returns all products
	GetAll(ctx context.Context) ([]*model.Product, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	err := db.AutoMigrate(&Product{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, product *model.Product) (uint, error) {
	mm := NewModel(product)

	if err := m.db.Create(mm).Error; err != nil {
		return 0, cmd.ParseMysqlError(ctx, "product", err)
	}
	return mm.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, product *model.Product) error {
	mm := NewModel(product)
	if err := m.db.Updates(mm).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "product", err)
	}
	return nil
}

func (m *MySqlRepository) Delete(ctx context.Context, id uint) error {
	if err := m.db.Delete(&Product{}, "id = ?", id).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "product", err)
	}
	return nil
}

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	var mm Product
	if err := m.db.First(&mm, "id = ?", id).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "product", err)
	}
	return mm.ToDto(), nil
}

func (m *MySqlRepository) GetAll(ctx context.Context) ([]*model.Product, error) {
	var models []*Product
	if err := m.db.Find(&models).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "product", err)
	}
	products := make([]*model.Product, len(models))
	for i, mm := range models {
		products[i] = mm.ToDto()
	}
	return products, nil
}
