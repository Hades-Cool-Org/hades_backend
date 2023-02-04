package user

import (
	"context"
	"gorm.io/gorm"
	"hades_backend/app/models/user"
)

type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user *user.User) error
	// Update updates an existing user
	Update(ctx context.Context, user *user.User) error
	// Delete deletes an existing user
	Delete(ctx context.Context, id string) error
	// GetByID returns a user by id
	GetByID(ctx context.Context, id string) (*user.User, error)
	// GetByEmail returns a user by email
	GetByEmail(ctx context.Context, email string) (*user.User, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	db.AutoMigrate(&User{}) //TODO

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, user *user.User) error {

	model := NewModel(user)

	return m.db.Create(model).Error
}

func (m *MySqlRepository) Update(ctx context.Context, user *user.User) error {
	return nil
}

func (m *MySqlRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MySqlRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	return nil, nil
}

func (m *MySqlRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u User

	result := m.db.Where("email = ?", email).First(u)

	if result.Error != nil {
		return nil, result.Error
	}

	return u.ToDto(), nil
}
