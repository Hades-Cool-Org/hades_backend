package store

import "gorm.io/gorm"

type Repository interface {
	// Create creates a new store
	Create(store *Store) (uint, error)
	// Update updates an existing store
	Update(store *Store) error
	// Delete deletes an existing store
	Delete(id uint) error
	// GetByID returns a store by id
	GetByID(id uint) (*Store, error)
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

func (m *MySqlRepository) Create(store *Store) (uint, error) {
	if err := m.db.Create(store).Error; err != nil {
		return 0, err
	}

	return store.ID, nil
}

func (m *MySqlRepository) Update(store *Store) error {
	return m.db.Updates(store).Error
}

func (m *MySqlRepository) Delete(id uint) error {
	return m.db.Delete(&Store{}, id).Error
}

func (m *MySqlRepository) GetByID(id uint) (*Store, error) {
	var store Store
	err := m.db.First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}
