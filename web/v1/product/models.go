package product

import (
	"errors"
	"net/http"
)

type Product struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Details       string `json:"details"`
	Image         string `json:"image_url"`
	MeasuringUnit string `json:"measuring_unit"`
}

type Request struct {
	*Product
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
	*Product
}

type GetAllResponse struct {
	Products []*Product `json:"products"`
}

func (p *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
