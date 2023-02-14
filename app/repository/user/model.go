package user

import (
	"database/sql"
	"gorm.io/gorm"
	"hades_backend/app/model/user"
)

type User struct {
	gorm.Model
	Name       string       `gorm:"type:varchar(255);not null;"`
	Email      string       `gorm:"type:varchar(255);not null;unique"`
	Phone      string       `gorm:"type:varchar(255);not null;"`
	Password   string       `gorm:"type:varchar(255);not null;"`
	FirstLogin sql.NullBool `gorm:"default:true"`
	Roles      []*Role      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) ToDto() *user.User {

	var roles []*user.Role

	for _, role := range u.Roles {
		roles = append(roles, &user.Role{Name: role.Name})
	}

	return &user.User{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		Roles:      roles,
		Password:   u.Password,
		FirstLogin: u.FirstLogin.Bool,
		DeletedAt:  u.DeletedAt.Time,
	}
}

type Role struct {
	ID     uint   `gorm:"primarykey"`
	UserId uint   `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE;"`
	Name   string `gorm:"type:varchar(255);not null"`
}

func NewModel(user *user.User) *User {

	var roles []*Role

	for _, role := range user.Roles {
		r := &Role{Name: role.Name}
		roles = append(roles, r)
	}

	u := &User{
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Roles:    roles,
	}

	if user.ID != 0 {
		u.ID = user.ID
	}

	return u
}
