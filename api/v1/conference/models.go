package conference

import (
	"errors"
	"hades_backend/app/model"
	"net/http"
)

type Response struct {
	*model.Occurrence
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Request struct {
	*model.Occurrence
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.StoreID == 0 {
		return errors.New("StoreID is required")
	}
	if r2.DeliveryID == 0 {
		return errors.New("DeliveryID is required")
	}

	if len(r2.Items) == 0 {
		return errors.New("items cannot be empty")
	}

	for _, product := range r2.Items {
		if product.ProductID == 0 {
			return errors.New("product id cannot be empty")
		}

		if product.Quantity == 0 {
			return errors.New("product expected quantity cannot be empty")
		}

	}

	return nil
}

type ListResponse struct {
	Occurrences []*model.Occurrence `json:"occurrences"`
}

func (c *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
