package auth

import (
	"errors"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"hades_backend/app/cmd/users"
	"time"
)

var (
	TokenAuth *jwtauth.JWTAuth
	ttl       = 4 * time.Hour
)

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

type Service struct {
}

func (u *Service) Login(email, password string) (*UserLogin, error) {

	logins := map[string]*users.UserDB{
		"oscar@gmail.com": {
			ID:    "id",
			Name:  "oscar",
			Email: "oscar@gmail.com",
			Phone: "",
			Roles: []*users.UserRolesDB{
				{
					ID:   "1",
					Name: "admin",
				},
			},
			Password: "password",
		},
		"guilherme@gmail.com": {
			ID:    "id",
			Name:  "guilherme",
			Email: "guilherme@gmail.com",
			Phone: "",
			Roles: []*users.UserRolesDB{
				{
					ID:   "1",
					Name: "admin",
				},
			},
			Password: "password",
		},
	}

	val, ok := logins[email]
	// If the key exists
	if ok {
		if val.Password == password {
			return &UserLogin{Token: encodeUserToken(val)}, nil
		}
	}

	return nil, errors.New("invalid user or password")
}

func encodeUserToken(user *users.UserDB) string {

	var roles []string

	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"name":    user.Name,
		"roles":   roles,
		"exp":     time.Now().UTC().Add(ttl).Unix(),
	})

	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)

	return tokenString
}
