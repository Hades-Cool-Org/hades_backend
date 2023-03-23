package model

type Stock struct { //NO UUID, WILL BE A SELECT ALL QUERY
	ID           uint           `json:"id"`
	StoreId      uint           `json:"store_id"`
	StoreName    string         `json:"store_name"`
	LastModified string         `json:"last_modified"`
	Products     []*ProductData `json:"products"`
}

type ProductData struct {
	ID          uint    `json:"id"`
	ProductId   uint    `json:"product_id"`
	ProductName string  `json:"name"`
	ImageUrl    string  `json:"image_url"`
	Current     float32 `json:"current"`
	Suggested   float32 `json:"suggested"`
}
