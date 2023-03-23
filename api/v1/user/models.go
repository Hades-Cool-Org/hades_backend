package user

import (
	"errors"
	"hades_backend/app/model"
	"net/http"
)

type Request struct {
	*model.User
}

type Response struct {
	*model.User
}

func (u *Response) Render(w http.ResponseWriter, r *http.Request) error {
	u.Password = "***"
	return nil
}

func (u *Request) Bind(r *http.Request) error {

	if u.Name == "" {
		return errors.New("id cannot be null")
	}

	//todo: additional validations

	return nil
}
