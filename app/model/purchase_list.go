package model

type PurchaseList struct {
	ID       uint       `json:"id"`
	User     *User      `json:"user"`
	Products []*Product `json:"products"`
}
