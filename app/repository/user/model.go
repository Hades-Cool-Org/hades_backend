package user

import (
	"database/sql"
	"gorm.io/gorm"
	"hades_backend/app/models/user"
	"strconv"
)

type User struct {
	gorm.Model
	Name       string       `gorm:"type:varchar(255);not null"`
	Email      string       `gorm:"type:varchar(255);not null"`
	Phone      string       `gorm:"type:varchar(255);not null"`
	Password   string       `gorm:"type:varchar(255);not null"`
	Created    int64        `gorm:"autoCreateTime:nano"`
	Updated    int64        `gorm:"autoUpdateTime:nano"`
	FirstLogin sql.NullBool `gorm:"default:true"`
	Roles      []*Role      `gorm:"foreignKey:UserId"`
}

func (u *User) ToDto() *user.User {

	var roles []*user.Role

	for _, role := range u.Roles {
		roles = append(roles, &user.Role{Name: role.Name})
	}

	return &user.User{
		ID:         "teste",
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
	gorm.Model
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

	if user.ID != "" {
		x, _ := strconv.Atoi(user.ID)
		u.ID = uint(x)
	}

	return u
}
