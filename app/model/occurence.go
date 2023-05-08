package model

type Occurrence struct {
	ID          uint              `json:"id"`
	DeliveryID  uint              `json:"delivery_id"`
	StoreID     uint              `json:"store_id"`
	User        *User             `json:"user"`
	Items       []*OccurrenceItem `json:"items"`
	CreatedDate string            `json:"created_date"`
	EndDate     string            `json:"end_date,omitempty"`
}

type OccurrenceType string

const (
	TypeCredit OccurrenceType = "CREDIT"
	TypeDebit  OccurrenceType = "DEBIT"
)

type OccurrenceItem struct {
	ProductID uint           `json:"product_id"`
	Type      OccurrenceType `json:"type"` //positive negative

	Name          string  `json:"name"`
	MeasuringUnit string  `json:"measuring_unit"`
	Quantity      float64 `json:"quantity"`
}
