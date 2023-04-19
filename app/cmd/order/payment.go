package order

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Type  string
	Total decimal.Decimal `gorm:"type:decimal(12,3);"`
	Text  string          `gorm:"type:text"`

	OrderID uint
}

func (p Payment) TableName() string {
	return "payments"
}
