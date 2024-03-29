package user

import (
	"context"
	"hades_backend/app/model"
	"testing"
)

func TestAuthService_Login(t *testing.T) {

	fl := true
	repo := &MockRepository{
		Users: map[uint]*model.User{
			1: {
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "cGFzc3dvcmQ=",
				Roles: []*model.Role{
					{Name: "user"},
				},
				FirstLogin: &fl,
			},
		},
	}

	authService := NewAuth(repo)

	ctx := context.Background()
	email := "test@example.com"
	password := "password"
	login, err := authService.Login(ctx, email, password)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if login == nil {
		t.Errorf("expected login to be non-nil")
	}
	if login.Token == "" {
		t.Errorf("expected token to be non-empty string")
	}
}
