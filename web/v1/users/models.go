package users

import (
	"errors"
	"net/http"
)

type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Phone string   `json:"phone"`
	Roles []*Roles `json:"roles"`
}

type Roles struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Request struct {
	*User
}

type Response struct {
	*User
}

func (u *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *Request) Bind(r *http.Request) error {

	if u.Name == "" {
		return errors.New("id cannot be null")
	}

	//todo: additional validations

	return nil
}
