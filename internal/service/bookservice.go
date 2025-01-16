package service

import "library/internal/domain/models"

type BookStorage interface {
	GetBooks() ([]models.BookStruct, error)
	GetBook(string) (models.BookStruct, error)
	SaveBook(models.BookStruct) (string, error)
	DeleteBook(string) error
}

type BookServiceStruct struct {
	storage BookStorage
}

func NewBookService(storage BookStorage) BookServiceStruct {

	return BookServiceStruct{storage: storage}

}

func (bs BookServiceStruct) GetBooks() ([]models.BookStruct, error) {
	return bs.storage.GetBooks()
}

func (bs BookServiceStruct) GetBook(id string) (models.BookStruct, error) {
	return bs.storage.GetBook(id)
}

func (bs BookServiceStruct) AddBook(book models.BookStruct) (string, error) {
	return bs.storage.SaveBook(book)
}

func (bs BookServiceStruct) DeleteBook(id string) error {
	return bs.storage.DeleteBook(id)
}
