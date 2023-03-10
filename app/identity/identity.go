package identity

import (
	"hades_backend/app/config"
)

type Identity struct {
	UserId uint
	Roles  []string
}

const (
	roleAdmin  = "admin"
	roleBuyer  = "buyer"
	roleDriver = "driver"
)

func (i *Identity) IsAdmin() bool {

	if !config.IsProd() {
		return true
	}

	for _, role := range i.Roles {
		if role == roleAdmin {
			return true
		}
	}
	return false
}

func (i *Identity) IsBuyer() bool {
	for _, role := range i.Roles {
		if role == roleBuyer {
			return true
		}
	}
	return false
}

func (i *Identity) IsDriver() bool {
	for _, role := range i.Roles {
		if role == roleDriver {
			return true
		}
	}
	return false
}
