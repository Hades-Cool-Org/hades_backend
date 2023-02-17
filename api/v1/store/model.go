package store

import (
	"errors"
	"net/http"

	"hades_backend/app/model/store"
)

type Request struct {
	*store.Store
}

type AddCourierRequest struct {
	Couriers []*store.User `json:"couriers"`
}

func (a *AddCourierRequest) Bind(r *http.Request) error {

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
	return nil
}

type Response struct {
	*store.Store
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetAllResponse struct {
	Stores []*store.Store `json:"stores"`
}

func (g *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
