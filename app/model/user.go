package model

import "time"

type User struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Roles      []*Role   `json:"roles,omitempty"`
	Password   string    `json:"password,omitempty"`
	FirstLogin bool      `json:"first_login,omitempty"`
	DeletedAt  time.Time `json:"-"`
}

type Role struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Login struct {
	Token      string
	FirstLogin bool `json:"first_login,omitempty"`
}
