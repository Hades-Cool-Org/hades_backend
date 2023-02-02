package order

import (
	"errors"
	"net/http"
)

type Order struct {
	ID        string     `json:"id"`
	Vendor    *Vendor    `json:"vendor"`
	StartDate string     `json:"start_date"`
	State     string     `json:"state"` //"CRIADO,ACEITO,ACEITO_PARCIAL,RECEBIDO,RECEBIDO_PARCIAL",
	EndDate   *string    `json:"end_date"`
	User      *User      `json:"user"`
	Total     string     `json:"total"`
	Payments  []*Payment `json:"payments"`
	Products  []*Product `json:"products"`
}

type Request struct {
	*Order
}

func (p *Request) Bind(r *http.Request) error {

	if p.Vendor == nil {
		return errors.New("vendor is requiredl")
	}

	if p.Vendor.ID == "" {
		return errors.New("vendorId is required")
	}

	if p.User == nil {
		return errors.New("user is required")
	}

	if p.User.ID == "" {
		return errors.New("user id is required")
	}
	return nil
}

type Response struct {
	*Order
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ListResponse struct {
	Orders []*Order `json:"orders"`
}

func (r3 *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// User comprador
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Vendor struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Payment struct {
	Type  string `json:"type"`
	Total string `json:"total"`
	Date  string `json:"date"`
}

type Product struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Image         string   `json:"image_url"`
	MeasuringUnit string   `json:"measuring_unit"`
	Quantity      float32  `json:"quantity"`
	Available     float32  `json:"available"` // quando houver uma coleta, alterar esse valor
	Total         string   `json:"total"`     //money TODO: RETORNAR UM VALOR INTEIRO?
	Stores        []*Store `json:"stores"`
}

type AddProductRequest struct {
	Products []*Product `json:"products"`
}

func (p *AddProductRequest) Bind(r *http.Request) error {

	if len(p.Products) == 0 {
		return errors.New("products cannot be empty")
	}

	for _, product := range p.Products {
		if product.ID == "" {
			return errors.New("product id is required")
		}

		errFun := func(message string) error {
			return errors.New("ProductId: " + product.ID + " -> " + message)
		}

		if product.Quantity == 0 {
			return errFun("quantity cannot be 0")
		}

		if product.Total == "" {
			return errFun("Total cannot be empty")
		}

		if len(product.Stores) == 0 {
			return errFun("store cannot be empty")
		}

		for _, store := range product.Stores {

			if store.ID == "" {
				return errFun("storeId cannot be empty")
			}

			if store.Quantity == 0 {
				return errFun("quantity cannot be 0")
			}
		}
	}

	return nil
}

type UpdateProductRequest struct {
	*Product
}

func (p *UpdateProductRequest) Bind(r *http.Request) error {

	errFun := func(message string) error {
		return errors.New("ProductId: " + p.ID + " -> " + message)
	}

	if p.Quantity == 0 {
		return errFun("quantity cannot be 0")
	}

	if len(p.Stores) == 0 {
		return errFun("store cannot be empty")
	}

	for _, store := range p.Stores {

		if store.ID == "" {
			return errFun("storeId cannot be empty")
		}

		if store.Quantity == 0 {
			return errFun("quantity cannot be 0")
		}
	}

	return nil
}

type ProductResponse struct {
	*Product
}

func (p ProductResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Store struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Quantity float32 `json:"quantity"`
}
