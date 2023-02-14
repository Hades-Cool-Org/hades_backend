package product

import (
	"errors"
	"hades_backend/app/model/product"
	"net/http"
)

type Request struct {
	*product.Product
}

func (p *Request) Bind(r *http.Request) error {
	if p.MeasuringUnit == "" {
		return errors.New("measuring unit cannot be null")
	}

	if p.Name == "" {
		return errors.New("name cannot be null")
	}

	return nil
}

type Response struct {
	*product.Product
}

type GetAllResponse struct {
	Products []*product.Product `json:"products"`
}

func (p *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
