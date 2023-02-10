package identity

type Identity struct {
	UserId uint
	Roles  []string
}

func (i *Identity) IsAdmin() bool {
	for _, role := range i.Roles {
		if role == "admin" {
			return true
		}
	}
	return false
}

func (i *Identity) IsBuyer() bool {
	for _, role := range i.Roles {
		if role == "buyer" {
			return true
		}
	}
	return false
}

func (i *Identity) IsDriver() bool {
	for _, role := range i.Roles {
		if role == "driver" {
			return true
		}
	}
	return false
}
