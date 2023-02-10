package user

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"hades_backend/app/model/user"
	"hades_backend/app/repository"
)

type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user *user.User) (uint, error)
	// Update updates an existing user
	Update(ctx context.Context, user *user.User) error
	// Delete deletes an existing user
	Delete(ctx context.Context, id uint) error
	// GetByID returns a user by id
	GetByID(ctx context.Context, id uint) (*user.User, error)
	// GetByEmail returns a user by email
	GetByEmail(ctx context.Context, email string) (*user.User, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	err := db.AutoMigrate(&User{}, &Role{})

	if err != nil {
		panic("oops!")
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, user *user.User) (uint, error) {
	model := NewModel(user)

	if err := m.db.Create(model).Error; err != nil {
		return 0, repository.ParseMysqlError("user", err)
	}

	return model.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, user *user.User) error {
	model := NewModel(user)
	return repository.ParseMysqlError("user",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&Role{}, "user_id = ?", user.ID).Error; err != nil {
				return err
			}

			if err := tx.Where("id = ?", user.ID).Updates(model).Error; err != nil {
				return err
			}

			return nil
		}),
	)
}

func (m *MySqlRepository) Delete(ctx context.Context, id uint) error {
	return repository.ParseMysqlError("user",
		m.db.Delete(&User{}, id).Error,
	)
}

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var u User

	if err := m.db.Where("id = ?", id).Preload("Roles").First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return u.ToDto(), nil
}

func (m *MySqlRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u User

	if err := m.db.Where("email = ?", email).Preload("Roles").First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return u.ToDto(), nil
}
