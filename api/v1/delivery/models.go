package delivery

import (
	"errors"
	"net/http"
)

type Request struct {
	*Delivery
}

type CompleteDeliveryRequest struct {
	ID       string     `json:"id"`
	Products []*Product `json:"products"`
}

func (r2 *CompleteDeliveryRequest) Bind(r *http.Request) error {

	if len(r2.Products) == 0 {
		return errors.New("products cannot be empty")
	}

	for _, product := range r2.Products {

		if product.ID == "" {
			return errors.New("product.UUID cannot be empty")
		}

		//todo deveria validar quantidade?
	}
	return nil
}

type Response struct {
	*Delivery
}

type CompleteDeliveryResponse struct {
	*CompleteDeliveryRequest
}

func (c *CompleteDeliveryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.Vehicle.ID == "" {
		errors.New("vehicle UUID cannot be empty")
	}

	if r2.Order.ID == "" {
		errors.New("order UUID cannot be empty")
	}

	if r2.User.ID == "" {
		errors.New("user UUID cannot be empty")
	}

	return nil
}

func (r2 *StartDeliveryRequest) Bind(r *http.Request) error {

	if r2.CourierID == "" {
		return errors.New("courier_id cannot be empty")
	}
	return nil
}

type ListResponse struct {
	Deliveries []*Delivery `json:"deliveries"`
}

func (l *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
