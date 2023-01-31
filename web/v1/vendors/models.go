package vendors

import (
	"errors"
	"net/http"
)

type Vendor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Cnpj     string `json:"cnpj"`
	Type     string `json:"type"`
	Location string `json:"location"`
	Contact  *User  `json:"contact"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Request struct {
	*Vendor
}

func (v *Request) Bind(r *http.Request) error {

	if v.Name == "" {
		return errors.New("name cannot be empty")
	}

	if v.Phone == "" {
		return errors.New("phone cannot be empty")
	}

	if v.Type == "" { //todo: validate types
		return errors.New("type cannot be empty")
	}

	if v.Location == "" {
		return errors.New("location cannot be empty")
	}

	return nil
}

type Response struct {
	*Vendor
}

type GetAllResponse struct {
	Vendors []*Vendor `json:"vendors"`
}

func (p *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
