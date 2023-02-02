package item_delivery

import (
	"errors"
	"net/http"
)

type Delivery struct {
	ID        string  `json:"id"`
	State     string  `json:"state"`    //ABERTO,COLETADO,ENTREGUE
	OrderID   string  `json:"order_id"` //optional
	User      *User   `json:"user"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	Items     []*Item `json:"items"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	Type     string  `json:"type"`
	Quantity float32 `json:"quantity"`
}

type Request struct {
	*Delivery
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.User.ID == "" {
		return errors.New("user ID cannot be empty")
	}

	if len(r2.Items) == 0 {
		return errors.New("boxes cannot be empty")
	}

	for _, box := range r2.Items {

		if box.Type == "" {
			return errors.New("box.Type cannot be empty")
		}

		if box.Quantity == 0 {
			return errors.New("box.Quantity cannot be empty")
		}
	}

	return nil
}

type Response struct {
	*Delivery
}

func (c *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ListResponse struct {
	Deliveries []*Delivery `json:"deliveries"`
}

func (c *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
