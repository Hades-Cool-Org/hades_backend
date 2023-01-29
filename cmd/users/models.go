package users

type UserDB struct {
	ID       string
	Name     string
	Email    string
	Phone    string
	Roles    []*UserRolesDB
	Password string
}

type UserRolesDB struct {
	ID   string
	Name string
}
