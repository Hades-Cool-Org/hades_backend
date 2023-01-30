package stock

import (
	"errors"
	"net/http"
)

type Product struct {
	ID        string  `json:"id"` //in model, need to add store_id
	Name      string  `json:"name"`
	Current   float32 `json:"current"`
	Suggested float32 `json:"suggested"`
}

type ProductRequest struct {
	*Product
}

func (i *ProductRequest) Bind(r *http.Request) error {

	if i.Current < 0 {
		return errors.New("current cannot be less than zero")
	}

	if i.Suggested < 0 {
		return errors.New("suggested cannot be less than zero")
	}

	return nil
}

type ProductResponse struct {
	*Product
}

func (i *ProductResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Stock struct { //NO ID, WILL BE A SELECT ALL QUERY
	StoreId      string     `json:"store_id"`
	LastModified string     `json:"last_modified"`
	Stock        []*Product `json:"stock"`
}

type Request struct {
	*Stock
}

func (r2 *Request) Bind(r *http.Request) error {
	//no validations
	return nil
}

type Response struct {
	*Stock
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
