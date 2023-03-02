package store

import (
	"gorm.io/gorm"
	"hades_backend/app/repository/user"
)

type Store struct {
	gorm.Model
	Name     string       `gorm:"type:varchar(255);not null;"`
	Address  string       `gorm:"type:varchar(255);not null;"`
	User     *user.User   `gorm:"many2many:store_owner;"`
	Couriers []*user.User `gorm:"many2many:store_couriers;"`
}
