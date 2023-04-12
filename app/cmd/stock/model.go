package stock

import (
	"gorm.io/gorm"
	"hades_backend/app/cmd/product"
	"hades_backend/app/cmd/store"
	"hades_backend/app/model"
)

type Stock struct {
	gorm.Model
	StoreID  uint
	Store    *store.Store
	Products []*ProductData `gorm:"foreignKey:StockID"`
}

// TableName overrides the table name used by ProductData to `stock_products`
func (Stock) TableName() string {
	return "stock"
}

type ProductData struct {
	Current   float32
	Suggested float32
	StockID   uint `gorm:"primaryKey;autoIncrement:false"`
	ProductID uint `gorm:"primaryKey;autoIncrement:false"`
	Product   *product.Product
}

// TableName overrides the table name used by ProductData to `stock_products`
func (ProductData) TableName() string {
	return "stock_products"
}

func NewModel(s *model.Stock) *Stock {
	var products []*ProductData

	fnProduct := func(p *model.ProductData) *ProductData {

		if p == nil {
			return nil
		}

		z := &ProductData{
			Current:   p.Current,
			Suggested: p.Suggested,
			ProductID: p.ProductId,
		}

		if s.ID != 0 {
			z.StockID = s.ID
		}

		return z
	}

	for _, data := range s.Products {
		products = append(products, fnProduct(data))
	}

	s2 := &Stock{
		StoreID:  s.StoreId,
		Products: products,
	}

	if s.ID != 0 {
		s2.ID = s.ID
	}

	return s2
}

func (s *Stock) ToDTO() *model.Stock {
	var products []*model.ProductData

	for _, data := range s.Products {
		products = append(products, data.ToDTO())
	}

	storeName := ""

	if s.Store != nil {
		storeName = s.Store.Name
	}

	s2 := &model.Stock{
		ID:           s.ID,
		StoreId:      s.StoreID,
		StoreName:    storeName,
		LastModified: s.UpdatedAt.String(),
		Products:     products,
	}

	return s2
}

func (p *ProductData) ToDTO() *model.ProductData {
	p2 := &model.ProductData{
		Current:   p.Current,
		Suggested: p.Suggested,
		ProductId: p.ProductID,
	}

	if p.Product != nil {
		p2.ProductName = p.Product.Name
		p2.ImageUrl = p.Product.ImageUrl
	}

	return p2
}
