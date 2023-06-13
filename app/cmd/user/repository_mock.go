package user

import (
	"context"
	"github.com/pkg/errors"
	"hades_backend/app/model"
)

type MockRepository struct {
	Users map[uint]*model.User
}

func (m *MockRepository) GetMultipleByIds(ctx context.Context, ids []uint) ([]*User, error) {

	return nil, nil
}

func (m *MockRepository) Create(ctx context.Context, user *model.User) (uint, error) {
	user.ID = uint(len(m.Users) + 1)
	m.Users[user.ID] = user
	return user.ID, nil
}

func (m *MockRepository) Update(ctx context.Context, user *model.User) error {
	if _, ok := m.Users[user.ID]; !ok {
		return errors.New("user not found")
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id uint) error {
	if _, ok := m.Users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.Users, id)
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	if user, ok := m.Users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockRepository) GetUsers(ctx context.Context) ([]*model.User, error) {
	var result []*model.User
	for _, user := range m.Users {
		result = append(result, user)
	}
	return result, nil
}
