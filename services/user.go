package services

import (
    "golang.org/x/crypto/bcrypt"
    "library-api/models"
    "library-api/repositories"
)

type UserService struct {
    Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
    return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(username, password string) (*models.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    user := &models.User{
        Username: username,
        Password: string(hashedPassword),
    }
    if err := s.Repo.Create(user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *UserService) Login(username, password string) (*models.User, error) {
    user, err := s.Repo.FindByUsername(username)
    if err != nil {
        return nil, err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, err
    }
    return user, nil
}