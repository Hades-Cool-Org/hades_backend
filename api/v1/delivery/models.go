package delivery

import (
	"errors"
	"hades_backend/app/model"
	"net/http"
)

type Request struct {
	*model.Delivery
}

type ItemRequest struct {
	DeliveryItems []*model.DeliveryItem `json:"items"`
}

func (r2 *ItemRequest) Bind(r *http.Request) error {

	if len(r2.DeliveryItems) == 0 {
		return errors.New("products cannot be empty")
	}

	for _, item := range r2.DeliveryItems {

		if item.ProductID == 0 {
			return errors.New("productID cannot be empty")
		}
		if item.StoreID == 0 {
			return errors.New("storeID cannot be empty")
		}

		if item.Quantity.IsZero() {
			return errors.New("quantity cannot be empty")
		}
	}
	return nil
}

type Response struct {
	*model.Delivery
}

type CompleteDeliveryResponse struct {
	*ItemRequest
}

func (c *CompleteDeliveryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.Session == nil {
		return errors.New("session cannot be empty")
	}

	if r2.Order == nil {
		return errors.New("order cannot be empty")
	}

	if r2.Session.ID == 0 {
		return errors.New("session ID cannot be empty")
	}

	if r2.Order.ID == 0 {
		return errors.New("order UUID cannot be empty")
	}

	if len(r2.DeliveryItems) == 0 {
		return errors.New("products cannot be empty")
	}

	for _, product := range r2.DeliveryItems {

		if product.ProductID == 0 {
			return errors.New("product cannot be empty")
		}

		if product.Quantity.IsZero() {
			return errors.New("quantity cannot be empty")
		}

		if product.StoreID == 0 {
			return errors.New("storeid cannot be empty")
		}
	}

	return nil
}

type SessionRequest struct {
	*model.Session
}

func (r2 *SessionRequest) Bind(r *http.Request) error {

	if r2.User.ID == 0 {
		return errors.New("courier_id cannot be empty")
	}

	if r2.Vehicle.ID == 0 {
		return errors.New("vehicle_id cannot be empty")
	}

	return nil
}

type SessionResponse struct {
	*model.Session
}

func (l *SessionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ListResponse struct {
	Deliveries []*model.Delivery `json:"deliveries"`
}

func (l *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type SessionListResponse struct {
	Sessions []*model.Session `json:"sessions"`
}

func (l *SessionListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type VehicleRequest struct {
	*model.Vehicle
}

func (r2 *VehicleRequest) Bind(r *http.Request) error {

	if r2.Type == "" {
		return errors.New("type cannot be empty")
	}

	if r2.Name == "" {
		return errors.New("name cannot be empty")
	}

	return nil
}

type VehicleResponse struct {
	*model.Vehicle
}

func (l *VehicleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ListVehicleResponse struct {
	Deliveries []*model.Vehicle `json:"vehicles"`
}

func (l *ListVehicleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
