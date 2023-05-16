package stock

import (
	"errors"
	"hades_backend/app/model"
	"net/http"
)

type Request struct {
	*model.Stock
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.Store.ID == 0 {
		return errors.New("storeId cannot be zero")
	}

	for _, item := range r2.Items {
		if item.Current.IsZero() {
			return errors.New("current cannot be less than zero")
		}

		if item.Suggested.IsNegative() {
			return errors.New("suggested cannot be less than zero")
		}

		if item.ProductID == 0 {
			return errors.New("ProductId cannot be zero")
		}

		return nil
	}

	//no validations
	return nil
}

type Response struct {
	*model.Stock
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
