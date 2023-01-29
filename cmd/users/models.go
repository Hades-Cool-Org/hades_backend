package users

type UserDB struct {
	ID         string
	Name       string
	Email      string
	Phone      string
	Roles      []*UserRolesDB
	Password   string
	FirstLogin bool
	Enabled    bool
}

type UserRolesDB struct {
	ID   string
	Name string
}
