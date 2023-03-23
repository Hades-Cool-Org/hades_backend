package purchase_list

import (
	"gorm.io/gorm"
	"hades_backend/app/cmd/product"
	"hades_backend/app/cmd/user"
	productModel "hades_backend/app/model"
)

type PurchaseList struct {
	gorm.Model
	UserID   uint
	User     *user.User
	Products []*product.Product `gorm:"many2many:purchase_list_products;"`
}

func NewModel(p *productModel.PurchaseList) *PurchaseList {
	var products []*product.Product

	fnProduct := func(p *productModel.Product) *product.Product {
		if p == nil || p.ID == 0 {
			return nil
		}
		z := &product.Product{}
		z.ID = p.ID
		return z
	}

	var u *user.User

	for _, pdt := range p.Products {
		products = append(products, fnProduct(pdt))
	}

	fnUser := func(u *productModel.User) *user.User {
		if u == nil || u.ID == 0 {
			return nil
		}
		z := &user.User{}
		z.ID = u.ID
		return z
	}

	if p.User != nil {
		u = fnUser(p.User)
	}

	p2 := &PurchaseList{
		User:     u,
		UserID:   u.ID,
		Products: products,
	}

	if p.ID != 0 {
		p2.ID = p.ID
	}

	return p2
}

func (p *PurchaseList) ToDTO() *productModel.PurchaseList {

	var products []*productModel.Product

	var u *productModel.User

	for _, pp := range p.Products {
		products = append(products, pp.ToDto())
	}

	if p.User != nil {
		u = p.User.ToDto()
	}

	p2 := &productModel.PurchaseList{
		ID:       u.ID,
		User:     u,
		Products: products,
	}

	if p.ID != 0 {
		p2.ID = p.ID
	}

	return p2
}
