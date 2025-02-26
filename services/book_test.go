package services

import (
	"errors"
	"testing"
	"library-api/cache"
	"library-api/models"
	"github.com/stretchr/testify/assert"
)

type mockBookRepo struct {
	createFunc func(book *models.Book) error
	findAllFunc func(limit, offset int) ([]models.Book, error)
}

func (m *mockBookRepo) Create(book *models.Book) error {
	return m.createFunc(book)
}

func (m *mockBookRepo) FindAll(limit, offset int) ([]models.Book, error) {
	return m.findAllFunc(limit, offset)
}

func (m *mockBookRepo) FindByID(id uint) (*models.Book, error) {
	return nil, nil // Stubbed
}

func (m *mockBookRepo) Update(book *models.Book) error {
	return nil // Stubbed
}

func (m *mockBookRepo) Delete(id uint) error {
	return nil // Stubbed
}

func TestBookService_CreateBook(t *testing.T) {
	repo := &mockBookRepo{
		createFunc: func(book *models.Book) error {
			return nil
		},
	}
	cache := cache.NewCache("localhost:6379") // Mock Redis (won’t connect)
	service := NewBookService(repo, cache)

	book, err := service.CreateBook("Test Book", "Author")
	assert.NoError(t, err)
	assert.Equal(t, "Test Book", book.Title)
	assert.Equal(t, "Author", book.Author)
}

func TestBookService_GetBooks_CacheMiss(t *testing.T) {
	repo := &mockBookRepo{
		findAllFunc: func(limit, offset int) ([]models.Book, error) {
			return []models.Book{{Title: "Book1", Author: "Author1"}}, nil
		},
	}
	cache := cache.NewCache("localhost:6379") // Mock Redis
	service := NewBookService(repo, cache)

	books, err := service.GetBooks(10, 0)
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Book1", books[0].Title)
}