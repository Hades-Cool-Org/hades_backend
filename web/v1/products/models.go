package products

import (
	"errors"
	"net/http"
)

type Product struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Details       string `json:"details"`
	Image         string `json:"image_url"`
	MeasuringUnit string `json:"measuring_unit"`
}

type ProductRequest struct {
	*Product
}

func (p *ProductRequest) Bind(r *http.Request) error {
	if p.MeasuringUnit == "" {
		return errors.New("measuring unit cannot be null")
	}

	if p.Name == "" {
		return errors.New("name cannot be null")
	}

	return nil
}

type ProductResponse struct {
	*Product
}

type AllProductResponse struct {
	Products []*Product `json:"products"`
}

func (p *AllProductResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *ProductResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
