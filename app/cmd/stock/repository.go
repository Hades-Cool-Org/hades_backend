package stock

import (
	"context"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
	"hades_backend/app/model"
)

type Repository interface {

	//Create creates a new stock
	Create(ctx context.Context, stock *model.Stock) (uint, error)

	//Update updates an existing stock
	Update(ctx context.Context, stock *model.Stock) error

	//Delete deletes an existing stock
	Delete(ctx context.Context, id uint) error

	//FindAll returns all stocks
	FindAll(ctx context.Context) ([]*model.Stock, error)

	//FindAllByStoreID returns all stocks by store id
	FindAllByStoreID(ctx context.Context, storeId uint) ([]*model.Stock, error)

	//FindByID returns a stock by id
	FindByID(ctx context.Context, id uint) (*model.Stock, error)

	//AddProductToStock adds a product to a stock
	AddProductToStock(ctx context.Context, stockId uint, products []*model.ProductData) error

	//RemoveProductFromStock removes a product from a stock
	RemoveProductFromStock(ctx context.Context, stockId uint, productId uint) error

	//UpdateProduct updates a product from a stock
	UpdateProduct(ctx context.Context, stockId uint, productId uint, product *model.ProductData) error

	//GetProduct returns a product from a stock
	GetProduct(ctx context.Context, stockId uint, productId uint) (*model.ProductData, error)
}

type MySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) Repository {

	err := db.AutoMigrate(&Stock{}, &ProductData{})

	if err != nil {
		panic(err)
	}

	return &MySQLRepository{db: db}
}

func (m *MySQLRepository) GetProduct(ctx context.Context, stockId uint, productId uint) (*model.ProductData, error) {
	var product ProductData

	err := cmd.ParseMysqlError(ctx, "stock",
		func() error {

			result := m.db.Where("stock_id = ? AND product_id = ?", stockId, productId).First(&product)

			if result.Error != nil {
				return result.Error
			}
			return nil

		}(),
	)

	if err != nil {
		return nil, err
	}

	return product.ToDTO(), err
}

func (m *MySQLRepository) UpdateProduct(ctx context.Context, stockId uint, productId uint, product *model.ProductData) error {

	return cmd.ParseMysqlError(ctx, "stock",
		func() error {

			stockResult, err := m.FindByID(ctx, stockId)

			if err != nil {
				return err
			}

			err = m.db.Model(&stockResult).
				Association("Items").
				Replace(&ProductData{ProductID: productId, StockID: stockId, Current: product.Current, Suggested: product.Suggested})

			if err != nil {
				return err
			}

			return nil
		}())

}

func (m *MySQLRepository) RemoveProductFromStock(ctx context.Context, stockId uint, productId uint) error {

	return cmd.ParseMysqlError(ctx, "stock",
		func() error {

			stockResult, err := m.FindByID(ctx, stockId)

			if err != nil {
				return err
			}

			err = m.db.Model(&stockResult).Association("Items").Delete(&ProductData{ProductID: productId, StockID: stockId})

			if err != nil {
				return err
			}

			return nil
		}())
}

func (m *MySQLRepository) Create(ctx context.Context, stock *model.Stock) (uint, error) {

	md := NewModel(stock)

	err := cmd.ParseMysqlError(ctx, "stock",
		m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Omit("Items").Create(md).Error; err != nil {
				return err
			}

			//err := tx.Model(md).Association("Items").Append(md.Items)
			//if err != nil {
			//	return err
			//}
			return nil
		}),
	)

	if err != nil {
		return 0, err
	}

	return md.ID, nil
}

func (m *MySQLRepository) AddProductToStock(ctx context.Context, stockId uint, products []*model.ProductData) error {

	return cmd.ParseMysqlError(ctx, "stock",
		func() error {

			stockResult, err := m.FindByID(ctx, stockId)

			if err != nil {
				return err
			}

			err = m.db.Model(&stockResult).Association("Items").Append(products)

			if err != nil {
				return err
			}

			return nil
		}())
}

func (m *MySQLRepository) FindAllByStoreID(ctx context.Context, storeId uint) ([]*model.Stock, error) {

	var models []*Stock

	err := cmd.ParseMysqlError(ctx, "stock",
		m.db.Where("store_id = ?", storeId).
			Preload("Items").
			Preload("Store").
			Preload("Items.Product").
			Find(&models).Error,
	)

	if err != nil {
		return nil, err
	}

	var returnModels []*model.Stock

	for _, model := range models {
		returnModels = append(returnModels, model.ToDTO())
	}

	return returnModels, nil
}

func (m *MySQLRepository) FindAll(ctx context.Context) ([]*model.Stock, error) {

	var models []*Stock

	err := cmd.ParseMysqlError(ctx, "stock",
		m.db.Find(&models).Error)

	if err != nil {
		return nil, err
	}

	var returnModels []*model.Stock

	for _, model := range models {
		returnModels = append(returnModels, model.ToDTO())
	}

	return returnModels, nil
}

func (m *MySQLRepository) FindByID(ctx context.Context, id uint) (*model.Stock, error) {

	var model Stock

	err := cmd.ParseMysqlError(ctx, "stock",
		m.db.Preload("Items").Preload("Store").Preload("Items.Product").First(&model, id).Error,
	)

	if err != nil {
		return nil, err
	}

	return model.ToDTO(), nil
}

func (m *MySQLRepository) Update(ctx context.Context, stock *model.Stock) error {

	model := NewModel(stock)

	return cmd.ParseMysqlError(ctx, "stock",
		m.db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Omit("Items").Save(model).Error; err != nil {
				return err
			}

			//err := tx.Model(model).Association("Items").Replace(model.Items)
			//if err != nil {
			//	return err
			//}
			return nil
		}),
	)
}

func (m *MySQLRepository) Delete(ctx context.Context, id uint) error {

	s := &Stock{}
	s.ID = id

	err := cmd.ParseMysqlError(ctx, "stock",
		m.db.Select("Items").Delete(&s).Error,
	)

	return err
}
