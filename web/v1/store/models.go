package store

import (
	"errors"
	"net/http"
)

type Store struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Request struct {
	*Store
}

func (r2 *Request) Bind(r *http.Request) error {

	if r2.Name == "" {
		return errors.New("name cannot be empty")
	}

	if r2.Address == "" {
		return errors.New("address cannot be empty")
	}

	return nil
}

type Response struct {
	*Store
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetAllResponse struct {
	Stores []*Store `json:"stores"`
}

func (g *GetAllResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
