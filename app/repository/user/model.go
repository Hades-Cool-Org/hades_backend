package user

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hades_backend/app/models/user"
)

type User struct {
	ID         uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()"` // primary key
	Name       string       `gorm:"not null"`
	Email      string       `gorm:"not null"`
	Phone      string       `gorm:"not null"`
	Password   string       `gorm:"not null"`
	Created    int64        `gorm:"autoCreateTime:nano"`
	Updated    int64        `gorm:"autoUpdateTime:nano"`
	FirstLogin sql.NullBool `gorm:"default:true"`
	DeletedAt  gorm.DeletedAt
	Roles      []*Role
}

func (u *User) ToDto() *user.User {

	var roles []*user.Role

	for _, role := range u.Roles {
		roles = append(roles, &user.Role{Name: role.Name})
	}

	return &user.User{
		ID:         u.ID.String(),
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
	Name string `gorm:"not null"`
}

func NewModel(user *user.User) *User {
	return &User{
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
	}
}
