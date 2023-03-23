package store

import (
	"context"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
	users "hades_backend/app/cmd/user"
	"hades_backend/app/model"
)

type Repository interface {
	// Create creates a new store
	Create(ctx context.Context, store *model.Store) (uint, error)
	// Update updates an existing store
	Update(ctx context.Context, store *model.Store) error
	// Delete deletes an existing store
	Delete(ctx context.Context, id uint) error
	// GetByID returns a store by id
	GetByID(ctx context.Context, id uint) (*model.Store, error)
	// GetAll returns all stores
	GetAll(ctx context.Context) ([]*model.Store, error)
	// GetByUserID returns all stores by user id
	GetByUserID(ctx context.Context, userId uint) ([]*model.Store, error)
	// RemoveCourierFromStore removes a courier from a store
	RemoveCourierFromStore(ctx context.Context, storeId uint, couriers []*users.User) error
}

type MySqlRepository struct {
	db *gorm.DB
}

func (m *MySqlRepository) RemoveCourierFromStore(ctx context.Context, storeId uint, couriers []*users.User) error {
	return cmd.ParseMysqlError(ctx, "store",
		func() error {
			storeResult, err := m.GetByID(ctx, storeId)

			if err != nil {
				return err
			}

			err = m.db.Model(&storeResult).Association("Couriers").Delete(couriers)

			if err != nil {
				return err
			}

			return nil
		}())
}

func (m *MySqlRepository) GetByUserID(ctx context.Context, userId uint) ([]*model.Store, error) {

	var models []*Store

	//TODO: not sure if we should keep it as raw query or use gorm
	err := m.db.Raw("SELECT * FROM stores where id in (SELECT store_id FROM hades_db.store_owner where user_id = ?)", userId).Scan(&models).Error

	if err != nil {
		return nil, err
	}

	var stores []*model.Store

	for _, model := range models {
		stores = append(stores, model.ToDTO())
	}

	return stores, nil
}

func (m *MySqlRepository) GetAll(ctx context.Context) ([]*model.Store, error) {
	var models []*Store

	err := m.db.Find(&models).Error

	if err != nil {
		return nil, err
	}

	var stores []*model.Store

	for _, mm := range models {
		stores = append(stores, mm.ToDTO())
	}

	return stores, nil
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	err := db.AutoMigrate(&Store{})

	if err != nil {
		panic("oops!")
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, store *model.Store) (uint, error) {
	mm := NewModel(store)

	err := cmd.ParseMysqlError(ctx, "store",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Omit("Couriers").Create(mm).Error; err != nil {
				return err
			}

			err := tx.Model(mm).Association("Couriers").Append(mm.Couriers)
			if err != nil {
				return err
			}
			return nil
		}),
	)

	if err != nil {
		return 0, err
	}

	return mm.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, store *model.Store) error {

	mm := NewModel(store)

	err := cmd.ParseMysqlError(ctx, "store",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("id = ?", mm.ID).Omit("Couriers").Updates(mm).Error; err != nil {
				return err
			}

			err := tx.Model(mm).Association("Couriers").Replace(mm.Couriers)

			if err != nil {
				return err
			}
			return nil
		}),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *MySqlRepository) Delete(ctx context.Context, id uint) error {

	s := &Store{}
	s.ID = id

	err := cmd.ParseMysqlError(ctx, "store",
		m.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&s).Association("User").Clear()
			if err != nil {
				return err
			}

			err = tx.Model(&s).Association("Couriers").Clear()
			if err != nil {
				return err
			}
			tx.Delete(&s)
			return nil
		}),
	)
	return err
}

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*model.Store, error) {
	var s Store
	err := cmd.ParseMysqlError(ctx, "store", m.db.Preload("Couriers").Preload("User").First(&s, id).Error)

	if err != nil {
		return nil, err
	}

	return s.ToDTO(), nil
}
