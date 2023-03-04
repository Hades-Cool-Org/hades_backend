package store

import (
	"gorm.io/gorm"
	"hades_backend/app/model/store"
	"hades_backend/app/repository/user"
)

type Store struct {
	gorm.Model
	Name     string       `gorm:"type:varchar(255);not null;"`
	Address  string       `gorm:"type:varchar(255);not null;"`
	User     *user.User   `gorm:"many2many:store_owner;"`
	Couriers []*user.User `gorm:"many2many:store_couriers;"`
}

func NewModel(s *store.Store) *Store {
	var couriers []*user.User

	fnUser := func(u *store.User) *user.User {
		if u == nil || u.ID == 0 {
			return nil
		}
		z := &user.User{}
		z.ID = u.ID
		return z
	}

	var u *user.User

	for _, courier := range s.Couriers {
		couriers = append(couriers, fnUser(courier))
	}

	if s.User != nil {
		u = fnUser(s.User)
	}

	return &Store{
		Name:     s.Name,
		Address:  s.Address,
		User:     u,
		Couriers: couriers,
	}
}

func (s *Store) ToDTO() *store.Store {

	var couriers []*store.User

	var u *store.User

	for _, courier := range s.Couriers {
		couriers = append(couriers, toStoreUser(courier))
	}

	if s.User != nil {
		u = toStoreUser(s.User)
	}

	return &store.Store{
		ID:       s.ID,
		Name:     s.Name,
		Address:  s.Address,
		User:     u,
		Couriers: couriers,
	}
}
func toStoreUser(u *user.User) *store.User {
	if u == nil {
		return nil
	}
	return &store.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}
