package services

import (
	"errors"
	"testing"
	"library-api/models"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	createFunc func(user *models.User) error
	findFunc   func(username string) (*models.User, error)
}

func (m *mockUserRepo) Create(user *models.User) error {
	return m.createFunc(user)
}

func (m *mockUserRepo) FindByUsername(username string) (*models.User, error) {
	return m.findFunc(username)
}

func TestUserService_CreateUser(t *testing.T) {
	repo := &mockUserRepo{
		createFunc: func(user *models.User) error {
			return nil
		},
	}
	service := NewUserService(repo)

	user, err := service.CreateUser("testuser", "password123")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.NotEmpty(t, user.Password) // Hashed
}

func TestUserService_Login_Success(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockUserRepo{
		findFunc: func(username string) (*models.User, error) {
			return &models.User{Username: "testuser", Password: string(hashed)}, nil
		},
	}
	service := NewUserService(repo)

	user, err := service.Login("testuser", "password123")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	repo := &mockUserRepo{
		findFunc: func(username string) (*models.User, error) {
			return &models.User{Username: "testuser", Password: "$2a$10$wronghash"}, nil
		},
	}
	service := NewUserService(repo)

	_, err := service.Login("testuser", "password123")
	assert.Error(t, err)
}