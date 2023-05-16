package model

import "github.com/shopspring/decimal"

type Stock struct { //NO UUID, WILL BE A SELECT ALL QUERY
	ID           uint         `json:"id"`
	Store        *Store       `json:"store"`
	LastModified string       `json:"last_modified"`
	Items        []*StockItem `json:"items"`
}

type StockItem struct {
	ProductID   uint            `json:"product_id"`
	ProductName string          `json:"name"`
	ImageUrl    string          `json:"image_url"`
	Current     decimal.Decimal `json:"current"`
	Suggested   decimal.Decimal `json:"suggested"`
	AvgPrice    decimal.Decimal `json:"avg_price"`
}
