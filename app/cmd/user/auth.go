package user

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
	"hades_backend/app/hades_errors"
	"hades_backend/app/logging"
	user2 "hades_backend/app/model/user"
	"hades_backend/app/repository/user"
	"time"
)

var (
	TokenAuth *jwtauth.JWTAuth
	ttl       = 4 * time.Hour
)

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

type AuthService struct {
	repository user.Repository
}

func NewAuth(r user.Repository) *AuthService {
	return &AuthService{repository: r}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*user2.Login, error) {

	logger := logging.FromContext(ctx)

	u, err := s.repository.GetByEmail(ctx, email)

	if err != nil {
		logger.Error("error getting u by email", zap.Error(err))
		return nil, err
	}

	if u == nil {
		return nil, hades_errors.NewForbiddenError(errors.New("invalid user or password"))
	}

	encodedPass := s.EncodePassword(password)

	println(encodedPass)

	if s.decodePassword(u.Password) == password {
		return &user2.Login{Token: s.encodeUserToken(u)}, nil
	}

	return nil, hades_errors.NewForbiddenError(errors.New("invalid user or password"))
}

func (s *AuthService) EncodePassword(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func (s *AuthService) decodePassword(password string) string {
	decoded, _ := base64.StdEncoding.DecodeString(password)
	return string(decoded)
}

func (s *AuthService) encodeUserToken(user *user2.User) string {

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

	return tokenString
}
