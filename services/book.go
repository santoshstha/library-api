package services

import (
	"context"
	// "encoding/json"
	"fmt"
	"library-api/cache"
	"library-api/models"
	"library-api/repositories"
	"time"
	"library-api/logger" // Import logger
)

type BookService struct {
	Repo  *repositories.BookRepository
	Cache *cache.Cache // Added Redis cache
}

func NewBookService(repo *repositories.BookRepository, cache *cache.Cache) *BookService {
	return &BookService{Repo: repo, Cache: cache}
}

func (s *BookService) CreateBook(title, author string) (*models.Book, error) {
	book := &models.Book{Title: title, Author: author}
	if err := s.Repo.Create(book); err != nil {
		return nil, err
	}
	// Invalidate cache on create (optional, if you want fresh data)
	s.Cache.Client.Del(context.Background(), "books")
	return book, nil
}

func (s *BookService) GetBooks(limit, offset int) ([]models.Book, error) {
	cacheKey := fmt.Sprintf("books:limit:%d:offset:%d", limit, offset)
	var books []models.Book

	err := s.Cache.Get(cacheKey, &books)
	if err == nil {
		logger.Logger.Log("Cache hit for books")
		return books, nil
	}

	books, err = s.Repo.FindAll(limit, offset)
	if err != nil {
		logger.Logger.Log(fmt.Sprintf("Error fetching books: %v", err))
		return nil, err
	}

	if err := s.Cache.Set(cacheKey, books, 5*time.Minute); err != nil {
		logger.Logger.Log(fmt.Sprintf("Failed to set cache: %v", err))
	}
	logger.Logger.Log("Books fetched from DB and cached")
	return books, nil
}

// Other methods (GetBook, UpdateBook, DeleteBook) remain unchanged
func (s *BookService) GetBook(id uint) (*models.Book, error) {
	return s.Repo.FindByID(id)
}

func (s *BookService) UpdateBook(id uint, title, author string) (*models.Book, error) {
	book, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	book.Title = title
	book.Author = author
	if err := s.Repo.Update(book); err != nil {
		return nil, err
	}
	// Invalidate books cache on update
	s.Cache.Client.Del(context.Background(), "books")
	return book, nil
}

func (s *BookService) DeleteBook(id uint) error {
	err := s.Repo.Delete(id)
	if err == nil {
		// Invalidate books cache on delete
		s.Cache.Client.Del(context.Background(), "books")
	}
	return err
}