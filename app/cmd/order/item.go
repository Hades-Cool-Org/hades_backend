package order

import (
	"github.com/shopspring/decimal"
	"hades_backend/app/cmd/product"
	"hades_backend/app/cmd/store"
)

type Item struct {
	OrderID uint `gorm:"primaryKey;autoIncrement:false"`

	ProductID uint `gorm:"primaryKey;autoIncrement:false"`
	Product   *product.Product

	StoreID uint `gorm:"primaryKey;autoIncrement:false"`
	Store   *store.Store

	Quantity  float64
	UnitPrice decimal.Decimal `gorm:"type:decimal(12,2);"`
}

func (i *Item) CalculateTotal() decimal.Decimal {
	return i.UnitPrice.Mul(decimal.NewFromFloat(i.Quantity))
}

func (i Item) TableName() string {
	return "order_items"
}
