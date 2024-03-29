package purchase_list

import (
	"errors"
	"hades_backend/app/model"
	"net/http"
)

type Request struct {
	*model.PurchaseList
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.User.ID == 0 {
		return errors.New("user id cannot be empty")
	}

	if len(r2.Products) == 0 {
		return errors.New("products cannot be empty")
	}

	for _, product := range r2.Products {
		if product.ID == 0 {
			return errors.New("product.id cannot be empty")
		}
	}

	return nil
}

type Response struct {
	*model.PurchaseList
}

type GetAllResponse struct {
	PurchaseLists []*model.PurchaseList `json:"lists"`
}

func (g *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (g *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
