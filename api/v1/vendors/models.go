package vendors

import (
	"errors"
	"hades_backend/app/model"
	"net/http"
)

type Request struct {
	*model.Vendor
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
	*model.Vendor
}

type GetAllResponse struct {
	Vendors []*model.Vendor `json:"vendors"`
}

func (p *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
