package purchase_list

import (
	"context"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
	"hades_backend/app/model/purchase_list"
)

type Repository interface {
	//Create creates a new purchase list
	Create(ctx context.Context, purchaseList *purchase_list.PurchaseList) (uint, error)
	//Update updates an existing purchase list
	Update(ctx context.Context, purchaseList *purchase_list.PurchaseList) error
	//Delete deletes an existing purchase list
	Delete(ctx context.Context, id uint) error
	//GetByID returns a purchase list by id
	GetByID(ctx context.Context, id uint) (*purchase_list.PurchaseList, error)
	//GetByUserID returns all purchase lists by user id
	GetByUserID(ctx context.Context, id uint) ([]*purchase_list.PurchaseList, error)
	//GetAll returns all purchase lists
	GetAll(ctx context.Context) ([]*purchase_list.PurchaseList, error)
}

type MySqlRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	// Migrate the schema
	err := db.AutoMigrate(&PurchaseList{})

	if err != nil {
		panic("oops!")
	}

	return &MySqlRepository{db: db}
}

func (m *MySqlRepository) Create(ctx context.Context, purchaseList *purchase_list.PurchaseList) (uint, error) {

	model := NewModel(purchaseList)

	err := cmd.ParseMysqlError(ctx, "purchase_list",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Omit("Items").Create(model).Error; err != nil {
				return err
			}

			err := tx.Model(model).Association("Items").Append(model.Products)
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

func (m *MySqlRepository) Update(ctx context.Context, purchaseList *purchase_list.PurchaseList) error {
	model := NewModel(purchaseList)

	err := cmd.ParseMysqlError(ctx, "purchase_list",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("id = ?", model.ID).Omit("Items").Updates(model).Error; err != nil {
				return err
			}

			err := tx.Model(model).Association("Items").Replace(model.Products)

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
	s := &PurchaseList{}
	s.ID = id

	err := cmd.ParseMysqlError(ctx, "store",
		m.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&s).Association("purchase_list").Clear()
			if err != nil {
				return err
			}

			err = tx.Model(&s).Association("Items").Clear()
			if err != nil {
				return err
			}
			tx.Delete(&s)
			return nil
		}),
	)
	return err
}

func (m *MySqlRepository) GetByID(ctx context.Context, id uint) (*purchase_list.PurchaseList, error) {
	s := &PurchaseList{}
	s.ID = id

	err := cmd.ParseMysqlError(ctx, "purchase_list",
		m.db.Preload("Items").Preload("User").First(s).Error,
	)

	if err != nil {
		return nil, err
	}

	return s.ToDTO(), nil
}

func (m *MySqlRepository) GetByUserID(ctx context.Context, id uint) ([]*purchase_list.PurchaseList, error) {
	s := &PurchaseList{}
	s.UserID = id

	var list []*PurchaseList
	err := cmd.ParseMysqlError(ctx, "purchase_list",
		m.db.Preload("Items").Preload("User").Where(s).Find(&list).Error,
	)

	if err != nil {
		return nil, err
	}

	var dtoList []*purchase_list.PurchaseList
	for _, item := range list {
		dtoList = append(dtoList, item.ToDTO())
	}

	return dtoList, nil
}

func (m *MySqlRepository) GetAll(ctx context.Context) ([]*purchase_list.PurchaseList, error) {
	var list []*PurchaseList
	err := cmd.ParseMysqlError(ctx, "purchase_list",
		m.db.Preload("Items").Find(&list).Error,
	)

	if err != nil {
		return nil, err
	}

	var dtoList []*purchase_list.PurchaseList
	for _, item := range list {
		dtoList = append(dtoList, item.ToDTO())
	}

	return dtoList, nil
}
