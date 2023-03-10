package purchase_list

import (
	"hades_backend/app/model/product"
	"hades_backend/app/model/user"
)

type PurchaseList struct {
	ID       uint               `json:"id"`
	User     *user.User         `json:"user"`
	Products []*product.Product `json:"products"`
}
