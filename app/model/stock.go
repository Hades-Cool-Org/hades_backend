package model

type Stock struct { //NO UUID, WILL BE A SELECT ALL QUERY
	ID           uint         `json:"id"`
	Store        *Store       `json:"store"`
	LastModified string       `json:"last_modified"`
	Items        []*StockItem `json:"items"`
}

type StockItem struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"name"`
	ImageUrl    string  `json:"image_url"`
	Current     float64 `json:"current"`
	Suggested   float64 `json:"suggested"`
}
