package users

import (
	"errors"
	"net/http"
)

type User struct {
	ID    string       `json:"id"`
	Name  string       `json:"name"`
	Email string       `json:"email"`
	Phone string       `json:"phone"`
	Roles []*UserRoles `json:"roles"`
}

type UserRoles struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserRequest struct {
	*User
}

type UserResponse struct {
	*User
}

func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *UserRequest) Bind(r *http.Request) error {

	if u.Name == "" {
		return errors.New("id cannot be null")
	}

	//todo: additional validations

	return nil
}
