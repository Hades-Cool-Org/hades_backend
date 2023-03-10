package stock

import (
	"errors"
	"hades_backend/app/model/stock"
	"net/http"
)

type ProductRequest struct {
	*stock.ProductData
}

type ProductRequestList struct {
	Products []*stock.ProductData `json:"products"`
}

func (i *ProductRequestList) Bind(r *http.Request) error {

	if len(i.Products) == 0 {
		return errors.New("products cannot be empty")
	}

	return nil
}

func (i *ProductRequest) Bind(r *http.Request) error {

	if i.Current < 0 {
		return errors.New("current cannot be less than zero")
	}

	if i.Suggested < 0 {
		return errors.New("suggested cannot be less than zero")
	}

	if i.ProductId == 0 {
		return errors.New("ProductId cannot be zero")
	}

	return nil
}

type ProductResponse struct {
	*stock.ProductData
}

func (i *ProductResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Request struct {
	*stock.Stock
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.StoreId == 0 {
		return errors.New("storeId cannot be zero")
	}

	//no validations
	return nil
}

type Response struct {
	*stock.Stock
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
