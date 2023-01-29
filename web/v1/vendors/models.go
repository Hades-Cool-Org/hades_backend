package vendors

import (
	"errors"
	"net/http"
)

type Vendor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Phone2   string `json:"phone2"`
	Cnpj     string `json:"cnpj"`
	Type     string `json:"type"`
	Location string `json:"location"`
}

type VendorRequest struct {
	*Vendor
}

func (v *VendorRequest) Bind(r *http.Request) error {

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

type VendorResponse struct {
	*Vendor
}

type AllVendorResponse struct {
	Vendors []*Vendor `json:"vendors"`
}

func (p *AllVendorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *VendorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
