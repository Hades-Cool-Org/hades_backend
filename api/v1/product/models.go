package product

import (
	"errors"
	"hades_backend/app/model"
	"net/http"
)

type Request struct {
	*model.Product
}

func (p *Request) Bind(r *http.Request) error {
	if p.MeasuringUnit == "" {
		return errors.New("measuring unit cannot be null")
	}

	if p.Name == "" {
		return errors.New("name cannot be null")
	}

	return nil
}

type Response struct {
	*model.Product
}

type GetAllResponse struct {
	Products []*model.Product `json:"products"`
}

func (p *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
