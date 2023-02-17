package store

type Store struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	User     *User   `json:"user"` //gerente da loja
	Couriers []*User `json:"couriers"`
}

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
