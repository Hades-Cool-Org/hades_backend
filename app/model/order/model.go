package order

import (
	"github.com/shopspring/decimal"
	"hades_backend/app/model/user"
	"hades_backend/app/model/vendors"
)

type Order struct {
	ID          uint            `json:"id"`
	Vendor      *vendors.Vendor `json:"vendor"`
	CreatedDate string          `json:"created_date"`
	State       string          `json:"state"` //"CRIADO,ACEITO,ACEITO_PARCIAL,RECEBIDO,RECEBIDO_PARCIAL",
	EndDate     *string         `json:"end_date"`
	User        *user.User      `json:"user"`
	Total       decimal.Decimal `json:"total"`

	Payments []*Payment `json:"payments"`
	Items    []*Item    `json:"items"`
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
	OrderID   uint `json:"order_id"`
	StoreID   uint `json:"store_id"`

	Name          string          `json:"name"`
	ImageUrl      string          `json:"image_url"`
	MeasuringUnit string          `json:"measuring_unit"`
	Quantity      float64         `json:"quantity"`
	Available     float64         `json:"available"` // quando houver uma coleta, alterar esse valor
	Total         decimal.Decimal `json:"total"`     //money TODO: RETORNAR UM VALOR INTEIRO?
}

func (i *Item) CalculateUnitPrice() decimal.Decimal {
	//TODO?

	if i.Total.IsZero() {
		return decimal.Zero
	}

	return i.Total.Div(decimal.NewFromFloat(i.Quantity))
}
