package login

import (
	"errors"
	"net/http"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lr *Request) Bind(r *http.Request) error {

	if lr.Password == "" {
		return errors.New("password cannot be null")
	}

	if lr.Email == "" {
		return errors.New("email cannot be null")
	}

	return nil
}

type Response struct {
	Token string `json:"jwt"`
}

func (u *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
