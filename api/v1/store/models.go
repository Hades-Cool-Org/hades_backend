package store

import (
	"errors"
	"net/http"
)

type Store struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	User     *User   `json:"user"` //gerente da loja
	Couriers []*User `json:"couriers"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Request struct {
	*Store
}

type AddCourierRequest struct {
	Couriers []*User `json:"couriers"`
}

func (a *AddCourierRequest) Bind(r *http.Request) error {

	if len(a.Couriers) == 0 {
		return errors.New("couriers cannot be empty")
	}

	for _, courier := range a.Couriers {
		if courier.ID == "" {
			return errors.New("courier.UUID cannot be empty")
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

	if r2.User.ID == "" {
		return errors.New("user cannot be empty")
	}
	return nil
}

type Response struct {
	*Store
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetAllResponse struct {
	Stores []*Store `json:"stores"`
}

func (g *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
