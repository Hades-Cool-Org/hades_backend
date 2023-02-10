package occurrence

import (
	"errors"
	"net/http"
)

type Occurrence struct {
	ID        string     `json:"id"`
	State     string     `json:"state"` //ABERTO,COLETADO,ENTREGUE
	OrderID   string     `json:"order_id"`
	Type      string     `json:"type"`  //positive negative
	Total     string     `json:"total"` //todo: mudar para numeric?
	User      *User      `json:"user"`
	Products  []*Product `json:"products"`
	StartDate string     `json:"start_date"`
	EndDate   string     `json:"end_date"`
}

type Product struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Image            string  `json:"image_url"`
	MeasuringUnit    string  `json:"measuring_unit"`
	Quantity         float32 `json:"quantity"`
	ExpectedQuantity float32 `json:"expected_quantity"`
}

type Response struct {
	*Occurrence
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Request struct {
	*Occurrence
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.OrderID == "" {
		return errors.New("orderID is required")
	}

	if r2.User == nil {
		return errors.New("user is required")
	}

	if r2.User.ID == "" {
		return errors.New("user id is required")
	}

	if len(r2.Products) == 0 {
		return errors.New("products cannot be empty")
	}

	for _, product := range r2.Products {
		if product.ID == "" {
			return errors.New("product id cannot be empty")
		}

		if product.ExpectedQuantity == 0 {
			return errors.New("product expected quantity cannot be empty")
		}

		if product.Quantity == 0 {
			return errors.New("product expected quantity cannot be empty")
		}
	}

	if r2.Total == "" {
		return errors.New("total cannot be empty")
	}

	return nil
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListResponse struct {
	Deliveries []*Occurrence `json:"occurrences"`
}

func (c *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
