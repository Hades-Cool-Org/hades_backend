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

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lr *UserLoginRequest) Bind(r *http.Request) error {

	if lr.Password == "" {
		return errors.New("password cannot be null")
	}

	if lr.Email == "" {
		return errors.New("email cannot be null")
	}

	return nil
}

type UserLoginResponse struct {
	Token string `json:"jwt"`
}

func (u *UserLoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
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
