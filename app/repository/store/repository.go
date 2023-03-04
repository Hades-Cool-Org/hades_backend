package store

import (
	"gorm.io/gorm"
	"hades_backend/app/model/store"
	"hades_backend/app/repository"
)

type Repository interface {
	// Create creates a new store
	Create(store *store.Store) (uint, error)
	// Update updates an existing store
	Update(store *store.Store) error
	// Delete deletes an existing store
	Delete(id uint) error
	// GetByID returns a store by id
	GetByID(id uint) (*store.Store, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewMySqlRepository(db *gorm.DB) *MySqlRepository {
	// Migrate the schema
	err := db.AutoMigrate(&Store{})

	if err != nil {
		panic("oops!")
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(store *store.Store) (uint, error) {

	model := NewModel(store)

	err := repository.ParseMysqlError("store",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Omit("User", "Couriers").Create(model).Error; err != nil {
				return err
			}

			err := tx.Model(model).Association("User").Replace(model.User)

			if err != nil {
				return err
			}

			err = tx.Model(model).Association("Couriers").Replace(model.Couriers)

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

func (m *MySqlRepository) Update(store *store.Store) error {
	_, err := m.Create(store)

	if err != nil {
		return err
	}

	return nil
}

func (m *MySqlRepository) Delete(id uint) error {

	s := &Store{}
	s.ID = id

	err := repository.ParseMysqlError("store",
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

func (m *MySqlRepository) GetByID(id uint) (*store.Store, error) {
	var s Store
	err := m.db.First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return s.ToDTO(), nil
}