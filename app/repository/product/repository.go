package product

import (
	"context"
	"gorm.io/gorm"
	"hades_backend/app/model/product"
	"hades_backend/app/repository"
)

type Repository interface {
	// Create creates a new product
	Create(ctx context.Context, product *product.Product) (uint, error)
	// Update updates an existing product
	Update(ctx context.Context, product *product.Product) error
	// Delete deletes an existing product
	Delete(ctx context.Context, id uint) error
	// GetByID returns a product by id
	GetByID(ctx context.Context, id uint) (*product.Product, error)
	// GetAll returns all products
	GetAll(ctx context.Context) ([]*product.Product, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	err := db.AutoMigrate(&Product{})

	if err != nil {
		panic("oops!")
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, product *product.Product) (uint, error) {
	model := NewModel(product)

	if err := m.db.Create(model).Error; err != nil {
		return 0, repository.ParseMysqlError(ctx, "product", err)
	}
	return model.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, product *product.Product) error {
	model := NewModel(product)
	if err := m.db.Updates(model).Error; err != nil {
		return repository.ParseMysqlError(ctx, "product", err)
	}
	return nil
}

func (m *MySqlRepository) Delete(ctx context.Context, id uint) error {
	if err := m.db.Delete(&Product{}, "id = ?", id).Error; err != nil {
		return repository.ParseMysqlError(ctx, "product", err)
	}
	return nil
}

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*product.Product, error) {
	var model Product
	if err := m.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, repository.ParseMysqlError(ctx, "product", err)
	}
	return model.ToDto(), nil
}

func (m *MySqlRepository) GetAll(ctx context.Context) ([]*product.Product, error) {
	var models []*Product
	if err := m.db.Find(&models).Error; err != nil {
		return nil, repository.ParseMysqlError(ctx, "product", err)
	}
	products := make([]*product.Product, len(models))
	for i, model := range models {
		products[i] = model.ToDto()
	}
	return products, nil
}
