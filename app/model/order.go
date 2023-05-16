package model

import (
	"fmt"
	"github.com/shopspring/decimal"
)

type OrderState string

const (
	Created           OrderState = "CRIADO"
	Accepted          OrderState = "COLETADO"
	AcceptedPartially OrderState = "COLETADO_PARCIAL"
	Received          OrderState = "RECEBIDO"
	ReceivedPartially OrderState = "RECEBIDO_PARCIAL"
	Completed         OrderState = "COMPLETADO"
)

var (
	mapOrderState = map[string]OrderState{
		"CRIADO":           Created,
		"COLETADO":         Accepted,
		"COLETADO_PARCIAL": AcceptedPartially,
		"RECEBIDO":         Received,
		"RECEBIDO_PARCIAL": ReceivedPartially,
		"COMPLETADO":       Completed,
	}
)

func (o OrderState) String() string {
	return string(o)
}

func OrderStateFromString(s string) (*OrderState, error) {
	if state, ok := mapOrderState[s]; ok {
		return &state, nil
	}
	return nil, fmt.Errorf("invalid order state: %s", s)
}

type Order struct {
	ID          uint            `json:"id"`
	Vendor      *Vendor         `json:"vendor"`
	CreatedDate string          `json:"created_date"`
	State       *OrderState     `json:"state"` //"CRIADO,ACEITO,ACEITO_PARCIAL,RECEBIDO,RECEBIDO_PARCIAL",
	EndDate     string          `json:"end_date,omitempty"`
	User        *User           `json:"user"`
	Total       decimal.Decimal `json:"total,omitempty"`

	Payments []*Payment   `json:"payments,omitempty"`
	Items    []*OrderItem `json:"items,omitempty"`
}

type Payment struct {
	ID    uint            `json:"id"`
	Type  string          `json:"type"`
	Total decimal.Decimal `json:"total"`
	Date  string          `json:"date"`
	Text  string          `json:"text"`
}

type OrderItem struct {
	ProductID uint `json:"product_id"`
	OrderID   uint `json:"order_id,omitempty"`
	StoreID   uint `json:"store_id,omitempty"`

	Name          string          `json:"name"`
	ImageUrl      string          `json:"image_url"`
	MeasuringUnit string          `json:"measuring_unit"`
	Quantity      decimal.Decimal `json:"quantity"`
	Total         decimal.Decimal `json:"total"` //money TODO: RETORNAR UM VALOR INTEIRO?
}

func (i *OrderItem) CalculateUnitPrice() decimal.Decimal {
	//TODO?

	if i.Total.IsZero() {
		return decimal.Zero
	}

	return i.Total.Div(i.Quantity).Round(2)
}
