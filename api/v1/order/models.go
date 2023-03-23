package order

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"hades_backend/app/model"
	"net/http"
)

type Request struct {
	*model.Order
}

type UpdateRequest struct {
	*model.Order
}

func (p *UpdateRequest) Bind(r *http.Request) error {

	for _, item := range p.Items {

		if item.ProductID == 0 {
			return errors.New("product id is required")
		}

		if item.StoreID == 0 {
			return errors.New("store id is required")
		}

		errFun := func(message string) error {
			return errors.New(fmt.Sprintf("ProductId: %d -> %s ", item.ProductID, message))
		}

		if item.Quantity == 0 {
			return errFun("quantity cannot be 0")
		}

		if item.Total.IsZero() {
			return errFun("Total cannot be empty")
		}

	}

	return nil
}

func (p *Request) Bind(r *http.Request) error {

	if p.Vendor == nil {
		return errors.New("vendor is requiredl")
	}

	if p.Vendor.ID == 0 {
		return errors.New("vendorId is required")
	}

	if p.User == nil {
		return errors.New("user is required")
	}

	if p.User.ID == 0 {
		return errors.New("user id is required")
	}

	for _, item := range p.Items {

		if item.ProductID == 0 {
			return errors.New("product id is required")
		}

		if item.StoreID == 0 {
			return errors.New("store id is required")
		}

		errFun := func(message string) error {
			return errors.New(fmt.Sprintf("ProductId: %d -> %s ", item.ProductID, message))
		}

		if item.Quantity == 0 {
			return errFun("quantity cannot be 0")
		}

		if item.Total == decimal.Zero {
			return errFun("Total cannot be empty")
		}

	}

	return nil
}

type Response struct {
	*model.Order
}

func (r2 *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ListResponse struct {
	Orders []*model.Order `json:"orders"`
}

func (r3 *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type DeleteItemsRequest struct {
	Items []*model.Item `json:"items"`
}

func (p *DeleteItemsRequest) Bind(r *http.Request) error {

	if len(p.Items) == 0 {
		return errors.New("products cannot be empty")
	}

	for _, item := range p.Items {

		if item.ProductID == 0 {
			return errors.New("product id is required")
		}

		if item.StoreID == 0 {
			return errors.New("store id is required")
		}

	}

	return nil
}

type UpdateItemRequest struct {
	*model.Item
}

func (p *UpdateItemRequest) Bind(r *http.Request) error {

	errFun := func(message string) error {
		return errors.New(fmt.Sprintf("ProductId: %d -> %s ", p.ProductID, message))
	}

	if p.Quantity == 0 {
		return errFun("quantity cannot be 0")
	}

	if p.Total == decimal.Zero {
		return errFun("Total cannot be empty")
	}

	return nil
}

type ItemResponse struct {
	*model.Item
}

func (p ItemResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ListItemResponse struct {
	Items []*model.Item `json:"items"`
}

func (p ListItemResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type PaymentRequest struct {
	*model.Payment
}

func (p *PaymentRequest) Bind(r *http.Request) error {

	if p.Type == "" {
		return errors.New("payment type is required")
	}

	if p.Total == decimal.Zero {
		return errors.New("payment total is required")
	}

	return nil
}
