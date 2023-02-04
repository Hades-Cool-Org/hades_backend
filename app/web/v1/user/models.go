package user

import (
	"errors"
	"hades_backend/app/models/user"
	"net/http"
)

type User struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Password   string   `json:"password"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Roles      []*Roles `json:"roles"`
	FirstLogin bool     `json:"first_login"`
}

func (u *User) ToModel() *user.User {

	var roles []*user.Role

	for _, role := range u.Roles {
		roles = append(roles, &user.Role{Name: role.Name})
	}

	return &user.User{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		Password:   u.Password,
		Roles:      roles,
		FirstLogin: u.FirstLogin,
	}

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
