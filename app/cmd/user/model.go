package user

import (
	"database/sql"
	"gorm.io/gorm"
	"hades_backend/app/model"
)

type User struct {
	gorm.Model
	Name       string       `gorm:"type:varchar(255);not null;"`
	Email      string       `gorm:"type:varchar(255);not null;unique"`
	Phone      string       `gorm:"type:varchar(255);not null;"`
	Password   string       `gorm:"type:varchar(255);not null;"`
	FirstLogin sql.NullBool `gorm:"default:true"`
	Roles      []*Role      `gorm:"many2many:user_roles;"`
}

func (u *User) ToDto() *model.User {

	var roles []*model.Role

	for _, role := range u.Roles {
		roles = append(roles, &model.Role{Name: role.Name})
	}

	return &model.User{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		Roles:      roles,
		Password:   u.Password,
		FirstLogin: &u.FirstLogin.Bool,
		DeletedAt:  u.DeletedAt.Time,
	}
}

type Role struct {
	Name string `gorm:"type:varchar(255);primary_key;"`
}

func NewModel(user *model.User) *User {

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

	if user.FirstLogin != nil {

		fl := *(user.FirstLogin)

		u.FirstLogin = sql.NullBool{Bool: fl, Valid: true}
	}

	if user.ID != 0 {
		u.ID = user.ID
	}

	return u
}
