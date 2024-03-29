package store

import (
	"gorm.io/gorm"
	"hades_backend/app/cmd/user"
	"hades_backend/app/model"
)

type Store struct {
	gorm.Model
	Name     string       `gorm:"type:varchar(255);not null;"`
	Address  string       `gorm:"type:varchar(255);not null;"`
	Couriers []*user.User `gorm:"many2many:store_couriers;"`
	UserID   uint
	User     *user.User
}

func NewModel(s *model.Store) *Store {
	var couriers []*user.User

	fnUser := func(u *model.User) *user.User {
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

	s2 := &Store{
		Name:     s.Name,
		Address:  s.Address,
		User:     u,
		UserID:   u.ID,
		Couriers: couriers,
	}

	if s.ID != 0 {
		s2.ID = s.ID
	}

	return s2
}

func (s *Store) ToDTO() *model.Store {

	var couriers []*model.User

	var u *model.User

	for _, courier := range s.Couriers {
		couriers = append(couriers, toStoreUser(courier))
	}

	if s.User != nil {
		u = toStoreUser(s.User)
	}

	return &model.Store{
		ID:       s.ID,
		Name:     s.Name,
		Address:  s.Address,
		User:     u,
		Couriers: couriers,
	}
}
func toStoreUser(u *user.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}
