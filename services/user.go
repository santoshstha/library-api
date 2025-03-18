package services

import (
	"fmt"
    "time"
	"golang.org/x/crypto/bcrypt"
	"library-api/email"
	"library-api/models"
	"library-api/repositories"
)

type UserService struct {
	Repo        *repositories.UserRepository
	EmailService *email.EmailService // New
}

func NewUserService(repo *repositories.UserRepository, emailService *email.EmailService) *UserService {
	return &UserService{Repo: repo, EmailService: emailService}
}

func (s *UserService) CreateUser(username, password, email string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email, // New
	}
	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}
	// Send welcome email
	// Send welcome email with a unique ID
	s.EmailService.Send(email, "Welcome to the Library!",
		fmt.Sprintf("Hi %s, welcome to our library system!", username),
		fmt.Sprintf("welcome_%s_%d", username, time.Now().Unix()))
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