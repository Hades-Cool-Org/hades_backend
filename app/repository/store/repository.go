package store

import (
	"context"
	"gorm.io/gorm"
	"hades_backend/app/model/store"
	"hades_backend/app/repository"
	users "hades_backend/app/repository/user"
)

type Repository interface {
	// Create creates a new store
	Create(ctx context.Context, store *store.Store) (uint, error)
	// Update updates an existing store
	Update(ctx context.Context, store *store.Store) error
	// Delete deletes an existing store
	Delete(ctx context.Context, id uint) error
	// GetByID returns a store by id
	GetByID(ctx context.Context, id uint) (*store.Store, error)
	// GetAll returns all stores
	GetAll(ctx context.Context) ([]*store.Store, error)
	// GetByUserID returns all stores by user id
	GetByUserID(ctx context.Context, userId uint) ([]*store.Store, error)
	// RemoveCourierFromStore removes a courier from a store
	RemoveCourierFromStore(ctx context.Context, storeId uint, couriers []*users.User) error
}

type MySqlRepository struct {
	db *gorm.DB
}

func (m *MySqlRepository) RemoveCourierFromStore(ctx context.Context, storeId uint, couriers []*users.User) error {
	return repository.ParseMysqlError(ctx, "store",
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

func (m *MySqlRepository) GetByUserID(ctx context.Context, userId uint) ([]*store.Store, error) {

	var models []*Store

	//TODO: not sure if we should keep it as raw query or use gorm
	err := m.db.Raw("SELECT * FROM stores where id in (SELECT store_id FROM hades_db.store_owner where user_id = ?)", userId).Scan(&models).Error

	if err != nil {
		return nil, err
	}

	var stores []*store.Store

	for _, model := range models {
		stores = append(stores, model.ToDTO())
	}

	return stores, nil
}

func (m *MySqlRepository) GetAll(ctx context.Context) ([]*store.Store, error) {
	var models []*Store

	err := m.db.Find(&models).Error

	if err != nil {
		return nil, err
	}

	var stores []*store.Store

	for _, model := range models {
		stores = append(stores, model.ToDTO())
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

func (m *MySqlRepository) Create(ctx context.Context, store *store.Store) (uint, error) {
	//l := logging.FromContext(ctx)

	model := NewModel(store)

	err := repository.ParseMysqlError(ctx, "store",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Omit("User", "Couriers").Create(model).Error; err != nil {
				return err
			}

			//TODO: FUCKING HATE GORM
			if err := tx.Exec("INSERT INTO store_owner (store_id, user_id) VALUES (?, ?)", model.ID, model.User.ID).Error; err != nil {
				return err
			}

			err := tx.Model(model).Association("Couriers").Replace(model.Couriers)
			if err != nil {
				return err
			}
			return nil
		}),
	)

	if err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (m *MySqlRepository) Update(ctx context.Context, store *store.Store) error {

	model := NewModel(store)

	err := repository.ParseMysqlError(ctx, "store",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("id = ?", model.ID).Omit("User", "Couriers").Updates(model).Error; err != nil {
				return err
			}

			if err := tx.Exec("UPDATE store_owner SET user_id = ? where store_id = ?", model.User.ID, model.ID).Error; err != nil {
				return err
			}

			err := tx.Model(model).Association("Couriers").Replace(model.Couriers)

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

	err := repository.ParseMysqlError(ctx, "store",
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

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*store.Store, error) {
	var s Store
	err := m.db.First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return s.ToDTO(), nil
}
