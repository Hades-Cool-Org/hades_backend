package user

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"hades_backend/app/models/user"
	"hades_backend/app/repository"
)

type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user *user.User) (uint, error)
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
	db.AutoMigrate(&User{}, &Role{}) //TODO

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, user *user.User) (uint, error) {

	model := NewModel(user)

	err := m.db.Create(model).Error

	if err != nil {
		return 0, repository.ParseMysqlError("user", err)
	}

	return model.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, user *user.User) error {
	model := NewModel(user)
	m.db.Delete(&Role{}, "user_id = ?", user.ID) //TODO
	m.db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", user.ID).Updates(model)
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

	result := m.db.Where("email = ?", email).Preload("Roles").First(&u)

	err := result.Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return u.ToDto(), nil
}
