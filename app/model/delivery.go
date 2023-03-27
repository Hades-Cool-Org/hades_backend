package model

type DeliveryState string

const (
	OPENED    DeliveryState = "ABERTO"
	COLLECTED DeliveryState = "COLETADO"
	DELIVERED DeliveryState = "ENTREGUE"
)

var (
	mappingDeliveryState = map[string]DeliveryState{
		"ABERTO":   OPENED,
		"COLETADO": COLLECTED,
		"ENTREGUE": DELIVERED,
	}
)

type Delivery struct {
	ID        uint           `json:"id"`
	State     *DeliveryState `json:"state"` //ABERTO,COLETADO,ENTREGUE
	StartDate string         `json:"start_date"`
	EndDate   *string        `json:"end_date"`

	Order   *Order   `json:"order"`
	Session *Session `json:"session"`

	DeliveryItems []*DeliveryItem `json:"items"`
}

type DeliveryItem struct {
	ProductID uint `json:"product_id"`
	StoreID   uint `json:"store_id"`

	Name          string  `json:"name"`
	ImageUrl      string  `json:"image_url"`
	MeasuringUnit string  `json:"measuring_unit"`
	Quantity      float64 `json:"quantity"`
}

type Vehicle struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Session struct {
	ID        uint     `json:"id"`
	User      *User    `json:"user"` //motorista
	Vehicle   *Vehicle `json:"vehicle"`
	StartDate string   `json:"start_date"`
	EndDate   *string  `json:"end_date,omitempty"`
}
