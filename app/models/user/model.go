package user

import "time"

type User struct {
	ID         uint
	Name       string
	Email      string
	Phone      string
	Roles      []*Role
	Password   string
	FirstLogin bool
	DeletedAt  time.Time
}

type Role struct {
	ID   uint
	Name string
}

type Login struct {
	Token string
}
