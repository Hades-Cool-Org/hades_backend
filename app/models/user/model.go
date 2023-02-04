package user

import "time"

type User struct {
	ID         string
	Name       string
	Email      string
	Phone      string
	Roles      []*Role
	Password   string
	FirstLogin bool
	DeletedAt  time.Time
}

type Role struct {
	Name string
}

type Login struct {
	Token string
}
