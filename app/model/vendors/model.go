package vendors

type Vendor struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Cnpj     string   `json:"cnpj"`
	Type     string   `json:"type"`
	Location string   `json:"location"`
	Contact  *Contact `json:"contact"`
}

type Contact struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}
