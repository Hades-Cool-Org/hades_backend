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

	Quantity  decimal.Decimal `gorm:"type:decimal(12,3);"`
	UnitPrice decimal.Decimal `gorm:"type:decimal(12,3);"`
}

func (i *Item) CalculateTotal() decimal.Decimal {
	return i.UnitPrice.Mul(i.Quantity).Round(2)
}

func (i Item) TableName() string {
	return "order_items"
}
