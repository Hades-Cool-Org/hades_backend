package store

import (
	"errors"
	"fmt"
	"hades_backend/app/model"
	"net/http"
)

type Request struct {
	*model.Store
}

type UpdateCouriersRequest struct {
	Couriers []*model.User `json:"couriers"`
}

func (a *UpdateCouriersRequest) Bind(r *http.Request) error {

	if len(a.Couriers) == 0 {
		return errors.New("couriers cannot be empty")
	}

	for _, courier := range a.Couriers {
		if courier.ID == 0 {
			return errors.New("courier.ID cannot be empty")
		}
	}

	return nil
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.Name == "" {
		return errors.New("name cannot be empty")
	}

	if r2.Address == "" {
		return errors.New("address cannot be empty")
	}

	if r2.User.ID == 0 {
		return errors.New("user cannot be empty")
	}

	for i, v := range r2.Couriers {
		if v.ID == 0 {
			return errors.New(fmt.Sprintf("courier at index [ %v ] cannot be empty", i))
		}
	}

	return nil
}

type Response struct {
	*model.Store
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetAllResponse struct {
	Stores []*model.Store `json:"stores"`
}

func (g *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
