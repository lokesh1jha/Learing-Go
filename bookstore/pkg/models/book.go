package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lokesh1jha/bookstore/pkg/config"
)

type Book struct {
	gorm.Model
	Name        string `gorm:"size:255;not null" json:"title"`
	Author      string `gorm:"size:255;not null" json:"author"`
	Publication string `gorm:"size:255;not null" json:"description"`
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{}) // migrate with empty struct
}

func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(id int64) (*Book, *gorm.DB) {
	var getBook Book
	db := db.Where("id = ?", id).Find(&getBook)
	return &getBook, db
}

func DeleteBook(id int64) Book {
	var book Book
	db.Where("id = ?", id).Delete(book)
	return book
}
