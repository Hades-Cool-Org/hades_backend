package product

import (
	"gorm.io/gorm"
	"hades_backend/app/model/product"
)

type Product struct {
	gorm.Model
	Name          string `gorm:"type:varchar(255);not null;"`
	Details       string `gorm:"type:varchar(255)"`
	Image         string `gorm:"type:varchar(255)"`
	MeasuringUnit string `gorm:"type:varchar(55)"`
}

func (p *Product) ToDto() *product.Product {
	return &product.Product{
		ID:            p.ID,
		Name:          p.Name,
		Details:       p.Details,
		Image:         p.Image,
		MeasuringUnit: p.MeasuringUnit,
	}
}

func NewModel(product *product.Product) *Product {
	p := &Product{
		Name:          product.Name,
		Details:       product.Details,
		Image:         product.Image,
		MeasuringUnit: product.MeasuringUnit,
	}

	if product.ID != 0 {
		p.ID = product.ID
	}

	return p
}
