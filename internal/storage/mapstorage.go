package storage

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"library/internal/domain/models"
	"library/internal/storage/storageerror"
	"time"
)

type MapStorage struct {
	userStorage map[string]models.UserStruct
	bookStorage map[string]models.BookStruct
}

func NewMapStorage() *MapStorage { // Откуда IDE знает что я хочу написать??? Она и эту строку сама сгенерировала

	return &MapStorage{userStorage: make(map[string]models.UserStruct),
		bookStorage: make(map[string]models.BookStruct)} /// И эту строку тоже

}

//type MapUserStorage struct {
//	userStorage map[string]models.UserStruct
//}
//
//func NewMapUserStorage() *MapUserStorage { // Откуда IDE знает что я хочу написать??? Она и эту строку сама сгенерировала
//
//	return &MapUserStorage{userStorage: make(map[string]models.UserStruct)}
//
//}
//
//type MapBookStorage struct {
//	bookStorage map[string]models.BookStruct
//}
//
//func NewMapBookStorage() *MapBookStorage {
//
//	return &MapBookStorage{bookStorage: make(map[string]models.BookStruct)}
//
//}

func (ms *MapStorage) GetUsers() ([]models.UserStruct, error) {

	if len(ms.userStorage) == 0 {
		return nil, storageerror.ErrUserStorageEmpty
	}

	var users []models.UserStruct

	for _, usr := range ms.userStorage {
		users = append(users, usr)
	}

	return users, nil

}

func (ms *MapStorage) GetUser(id string) (models.UserStruct, error) {

	user, ok := ms.userStorage[id]

	if !ok {
		return models.UserStruct{}, storageerror.ErrUserNotFound
	}

	return user, nil

}

func (ms *MapStorage) SaveUser(user models.UserStruct) (string, error) {

	for _, usr := range ms.userStorage {

		if usr.Email == user.Email {
			return "", storageerror.ErrUserAlreadyExist
		}

	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	user.Password = string(hash)

	id := uuid.New()
	idStr := id.String()

	user.ID = id
	user.DateRegistration = time.Now()

	ms.userStorage[idStr] = user

	return idStr, nil

}

func (ms *MapStorage) ValidateUser(user models.UserLoginStruct) (string, error) {

	for key, userMS := range ms.userStorage {

		if userMS.Email == user.Email {

			if err := bcrypt.CompareHashAndPassword([]byte(userMS.Password), []byte(user.Password)); err != nil {
				return "", storageerror.ErrUserInvalidPassword
			}

			return key, nil

		}

	}

	return "", storageerror.ErrUserNotFound

}

func (ms *MapStorage) EditUser(id string, user models.UserStruct) error {

	userMS, ok := ms.userStorage[id]

	if !ok {
		return storageerror.ErrUserNotFound
	}

	if userMS.Email != user.Email {

		for _, userCheck := range ms.userStorage {

			if userCheck.Email == user.Email {
				return fmt.Errorf("user with email %s already exists", user.Email)
			}

		}

	}

	if errCompare := bcrypt.CompareHashAndPassword([]byte(userMS.Password), []byte(user.Password)); errCompare != nil {

		hash, errGen := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

		if errGen != nil {
			return errGen
		}

		user.Password = string(hash)

	} else {
		user.Password = userMS.Password
	}

	user.ID = userMS.ID
	user.DateRegistration = userMS.DateRegistration

	ms.userStorage[id] = user

	return nil

}

func (ms *MapStorage) DeleteUser(id string) error {

	_, ok := ms.userStorage[id]

	if !ok {
		return storageerror.ErrUserNotFound
	}

	delete(ms.userStorage, id)

	return nil

}

func (ms *MapStorage) GetBooks() ([]models.BookStruct, error) {

	if len(ms.bookStorage) == 0 {
		return nil, storageerror.ErrBookStorageEmpty
	}

	var books []models.BookStruct

	for _, bk := range ms.bookStorage {
		books = append(books, bk)
	}

	return books, nil

}

func (ms *MapStorage) GetBook(id string) (models.BookStruct, error) {

	book, ok := ms.bookStorage[id] // IDE сама

	if !ok { // IDE сама
		return models.BookStruct{}, storageerror.ErrBookNotFound // IDE сама
	}

	return book, nil

}

func (ms *MapStorage) SaveBook(book models.BookStruct) (string, error) {

	for _, bk := range ms.bookStorage {

		if bk.Name == book.Name && bk.Author == book.Author {
			return "", storageerror.ErrBookAlreadyExist // IDE сама написала строку
		}

	}

	ID := uuid.New()
	IDStr := ID.String()

	book.ID = ID

	ms.bookStorage[IDStr] = book

	return IDStr, nil

}

func (ms *MapStorage) EditBook(id string, book models.BookStruct) error {

	bookMS, ok := ms.bookStorage[id]

	if !ok {
		return storageerror.ErrBookNotFound
	}

	if bookMS.Name != book.Name || bookMS.Author != book.Author {

		for _, bookCheck := range ms.bookStorage {

			if bookCheck.Name == book.Name && bookCheck.Author == book.Author {
				return fmt.Errorf("book with name %s and author %s already exists", book.Name, book.Author)
			}

		}

	}

	book.ID = bookMS.ID

	ms.bookStorage[id] = book

	return nil

}

func (ms *MapStorage) DeleteBook(id string) error {

	_, ok := ms.bookStorage[id]

	if !ok {
		return storageerror.ErrBookNotFound
	}

	delete(ms.bookStorage, id)

	return nil

}
func (ms *MapStorage) DeleteBooks() error {
	return nil
}
