package delivery

import (
	"errors"
	"net/http"
)

type Delivery struct {
	ID        string   `json:"id"`
	State     string   `json:"state"` //ABERTO,COLETADO,ENTREGUE
	Order     *Order   `json:"order"`
	Store     *Store   `json:"store"`
	User      *User    `json:"user"` //motorista
	Vehicle   *Vehicle `json:"vehicle"`
	StartDate string   `json:"start_date"`
	EndDate   *string  `json:"end_date"`
}

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

type Store struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	User    *User  `json:"user"` //gerente
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Vehicle struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID        string     `json:"id"`
	StartDate string     `json:"start_date"`
	User      *User      `json:"user"` //comprador
	Products  []*Product `json:"products"`
}

type Product struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Image         string  `json:"image_url"`
	MeasuringUnit string  `json:"measuring_unit"`
	Quantity      float32 `json:"quantity"`
	Total         string  `json:"total"` //money TODO: RETORNAR UM VALOR INTEIRO?
}

type StartDeliveryRequest struct {
	UserID    string  `json:"user_id"`
	CourierID string  `json:"courier_id"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date"`
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
