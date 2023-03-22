package user

import "time"

type User struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Roles      []*Role   `json:"roles"`
	Password   string    `json:"password"`
	FirstLogin bool      `json:"first_login"`
	DeletedAt  time.Time `json:"-"`
}

type Role struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Login struct {
	Token string
}
