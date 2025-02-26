package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
}