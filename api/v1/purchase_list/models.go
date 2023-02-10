package purchase_list

import (
	"errors"
	"net/http"
)

type List struct {
	ID       string     `json:"id"`
	UserID   string     `json:"user_id"`
	Products []*Product `json:"products"`
}

type Product struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Image         string `json:"image_url"`
	MeasuringUnit string `json:"measuring_unit"`
}

type Request struct {
	*List
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.UserID == "" {
		return errors.New("user UUID cannot be empty")
	}

	if len(r2.Products) == 0 {
		return errors.New("products cannot be empty")
	}

	for _, product := range r2.Products {
		if product.ID == "" {
			return errors.New("product.UUID cannot be empty")
		}
	}

	return nil
}

type Response struct {
	*List
}

type GetAllResponse struct {
	PurchaseLists []*List `json:"lists"`
}

func (g *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (g *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
