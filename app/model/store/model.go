package store

import "encoding/json"

type Store struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	User     *User   `json:"user"` //gerente da loja
	Couriers []*User `json:"couriers"`
}

func (s *Store) ToLoggableString() string {
	j, err := json.Marshal(s)

	if err != nil {
		return "{}"
	}
	return string(j)
}

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
