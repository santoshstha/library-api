package repositories

import (
    "gorm.io/gorm"
    "library-api/models"
	"library-api/logger"
	
)

var Logger *logger.AsyncLogger

type BookRepository struct {
    DB *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
    return &BookRepository{DB: db}
}

func (r *BookRepository) Create(book *models.Book) error {
    return r.DB.Create(book).Error
}

func (r *BookRepository) FindAll(limit, offset int) ([]models.Book, error) {
    var books []models.Book
    err := r.DB.Limit(limit).Offset(offset).Find(&books).Error
    return books, err
}

func (r *BookRepository) FindByID(id uint) (*models.Book, error) {
    var book models.Book
    err := r.DB.First(&book, id).Error
    return &book, err
}

func (r *BookRepository) Update(book *models.Book) error {
    return r.DB.Save(book).Error
}

func (r *BookRepository) Delete(id uint) error {
    return r.DB.Delete(&models.Book{}, id).Error
}