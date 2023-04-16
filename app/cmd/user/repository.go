package user

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
	// Create creates a new user
	Create(ctx context.Context, user *model.User) (uint, error)
	// Update updates an existing user
	Update(ctx context.Context, user *model.User) error
	// Delete deletes an existing user
	Delete(ctx context.Context, id uint) error
	// GetByID returns a user by id
	GetByID(ctx context.Context, id uint) (*model.User, error)
	// GetByEmail returns a user by email
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	// GetMultipleByIds returns multiple users by ids
	GetMultipleByIds(ctx context.Context, ids []uint) ([]*User, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	err := db.AutoMigrate(&User{}, &Role{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) GetMultipleByIds(ctx context.Context, ids []uint) ([]*User, error) {

	var models []*User

	err := m.db.Model(&User{}).Where("id IN ?", ids).Find(&models).Error // not preloading rules, should we?

	if err != nil {
		return nil, cmd.ParseMysqlError(ctx, "user", err)
	}

	return models, nil
}

func (m *MySqlRepository) Create(ctx context.Context, user *model.User) (uint, error) {
	mm := NewModel(user)

	if err := m.db.Create(mm).Error; err != nil {
		return 0, cmd.ParseMysqlError(ctx, "user", err)
	}

	return mm.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, user *model.User) error {
	mm := NewModel(user)
	return cmd.ParseMysqlError(ctx, "user",
		m.db.Transaction(func(tx *gorm.DB) error {
			tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&mm)
			err := tx.Model(&mm).Association("Roles").Replace(&mm.Roles)
			if err != nil {
				return err
			}
			return nil
		}),
	)
}

func (m *MySqlRepository) Delete(ctx context.Context, id uint) error {
	return cmd.ParseMysqlError(ctx, "user",
		m.db.Select("Roles").Unscoped().Delete(&User{}, id).Error,
	)
}

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var u User
	if err := m.db.Where("id = ?", id).Preload("Roles").First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return u.ToDto(), nil
}

func (m *MySqlRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u User

	if err := m.db.Where("email = ?", email).Preload("Roles").First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return u.ToDto(), nil
}
