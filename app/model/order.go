package model

import (
	"fmt"
	"github.com/shopspring/decimal"
)

type OrderState string

const (
	Created           OrderState = "CRIADO"
	Accepted          OrderState = "ACEITO"
	AcceptedPartially OrderState = "ACEITO_PARCIAL"
	Received          OrderState = "RECEBIDO"
	ReceivedPartially OrderState = "RECEBIDO_PARCIAL"
	Completed         OrderState = "COMPLETADO"
)

var (
	mapOrderState = map[string]OrderState{
		"CRIADO":           Created,
		"ACEITO":           Accepted,
		"ACEITO_PARCIAL":   AcceptedPartially,
		"RECEBIDO":         Received,
		"RECEBIDO_PARCIAL": ReceivedPartially,
		"COMPLETADO":       Completed,
	}
)

func OrderStateFromString(s string) (OrderState, error) {
	if state, ok := mapOrderState[s]; ok {
		return state, nil
	}
	return "", fmt.Errorf("invalid order state: %s", s)
}

type Order struct {
	ID          uint            `json:"id"`
	Vendor      *Vendor         `json:"vendor"`
	CreatedDate string          `json:"created_date"`
	State       OrderState      `json:"state"` //"CRIADO,ACEITO,ACEITO_PARCIAL,RECEBIDO,RECEBIDO_PARCIAL",
	EndDate     string          `json:"end_date,omitempty"`
	User        *User           `json:"user"`
	Total       decimal.Decimal `json:"total,omitempty"`

	Payments []*Payment `json:"payments,omitempty"`
	Items    []*Item    `json:"items,omitempty"`
}

type Payment struct {
	ID    uint            `json:"id"`
	Type  string          `json:"type"`
	Total decimal.Decimal `json:"total"`
	Date  string          `json:"date"`
	Text  string          `json:"text"`
}

type Item struct {
	ProductID uint `json:"product_id"`
	OrderID   uint `json:"order_id,omitempty"`
	StoreID   uint `json:"store_id,omitempty"`

	Name          string          `json:"name"`
	ImageUrl      string          `json:"image_url"`
	MeasuringUnit string          `json:"measuring_unit"`
	Quantity      float64         `json:"quantity"`
	Available     float64         `json:"available,omitempty"` // quando houver uma coleta, alterar esse valor
	Total         decimal.Decimal `json:"total"`               //money TODO: RETORNAR UM VALOR INTEIRO?
}

func (i *Item) CalculateUnitPrice() decimal.Decimal {
	//TODO?

	if i.Total.IsZero() {
		return decimal.Zero
	}

	return i.Total.Div(decimal.NewFromFloat(i.Quantity))
}
