package order

import (
	"github.com/shopspring/decimal"
	"hades_backend/app/model/user"
	"hades_backend/app/model/vendors"
)

type Order struct {
	ID        string          `json:"id"`
	Vendor    *vendors.Vendor `json:"vendor"`
	StartDate string          `json:"start_date"`
	State     string          `json:"state"` //"CRIADO,ACEITO,ACEITO_PARCIAL,RECEBIDO,RECEBIDO_PARCIAL",
	EndDate   *string         `json:"end_date"`
	User      *user.User      `json:"user"`
	Total     uint64          `json:"total"`
	Payments  []*Payment      `json:"payments"`
	Items     []*Item         `json:"products"`
}

type Payment struct {
	Type  string          `json:"type"`
	Total decimal.Decimal `json:"total"`
	Date  string          `json:"date"`
	Text  string          `json:"text"`
}

type Item struct {
	ProductID     uint            `json:"product_id"`
	Name          string          `json:"name"`
	Image         string          `json:"image_url"`
	MeasuringUnit string          `json:"measuring_unit"`
	Quantity      float64         `json:"quantity"`
	Available     float64         `json:"available"` // quando houver uma coleta, alterar esse valor
	Total         decimal.Decimal `json:"total"`     //money TODO: RETORNAR UM VALOR INTEIRO?
	StoreID       uint            `json:"store_id"`
}

func (i *Item) CalculateUnitPrice() decimal.Decimal {
	//TODO?
	return i.Total.Div(decimal.NewFromFloat(i.Quantity))
}
