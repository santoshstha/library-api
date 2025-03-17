package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"uniqueIndex"` // Indexed for fast lookup
	Password string `json:"password"`
	Email string  `json:"email"` 
}

type Book struct {
	gorm.Model
	Title  string `json:"title" gorm:"index"`  // Indexed for searches
	Author string `json:"author" gorm:"index"` // Indexed for searches
}