package storage

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"library/internal/domain/models"
	"library/internal/logger"
	"library/internal/storage/storageerror"
)

type MapStorage struct {
	userStorage map[string]models.UserStruct
	bookStorage map[string]models.BookStruct
}

func NewMapStorage() *MapStorage { // Откуда IDE знает что я хочу написать??? Она и эту строку сама сгенерировала

	return &MapStorage{userStorage: make(map[string]models.UserStruct),
		bookStorage: make(map[string]models.BookStruct)} /// И эту строку тоже

}

func (ms *MapStorage) SaveUser(user models.UserStruct) (string, error) {

	log := logger.Get()

	for _, usr := range ms.userStorage {

		if usr.Email == user.Email { // ??? IDE тоже сама подставила. Откуда она всё знает?
			//return errors.New("User already exists") // ??? IDE тоже сама подставила. Откуда она всё знает?
			return "", fmt.Errorf("user with email %s already exists", user.Email) // ??? IDE тоже сама подставила строку в кавычках
			// ??? Почему на уроке используем fmt
			// а не строку выше?
		}

	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost) // ??? Опять IDE сама всю строку сделала. Как?

	if err != nil {
		return "", err
	}

	user.Password = string(hash)

	ID := uuid.New()
	IDStr := ID.String()

	user.ID = ID

	ms.userStorage[IDStr] = user

	log.Debug().Any("user storage", ms.userStorage).Msg("Check user storage")

	return IDStr, nil

}

func (ms *MapStorage) ValidateUser(user models.UserLoginStruct) (string, error) {

	for key, usr := range ms.userStorage {

		if usr.Email == user.Email {

			if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(user.Password)); err != nil { // IDE сгенерировала строку
				return "", fmt.Errorf("invalid user password")
			}

			return key, nil

		}

	}

	return "", fmt.Errorf("user not found")

}

func (ms *MapStorage) GetBooks() ([]models.BookStruct, error) {

	if len(ms.bookStorage) == 0 {
		return nil, storageerror.ErrBookStorageEmpty
	}

	var books []models.BookStruct

	for _, bk := range ms.bookStorage { // IDE всё знает!?
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

	log := logger.Get()

	for _, bk := range ms.bookStorage {

		if bk.Name == book.Name && bk.Author == book.Author {
			return "", storageerror.ErrBookAlreadyExist // IDE сама написала строку
		}

	}

	ID := uuid.New()
	IDStr := ID.String()

	book.ID = ID

	ms.bookStorage[IDStr] = book

	log.Debug().Any("book storage", ms.bookStorage).Msg("Check book storage")

	return IDStr, nil

}

func (ms *MapStorage) DeleteBook(id string) error {

	_, ok := ms.bookStorage[id]

	if !ok {
		return storageerror.ErrBookNotFound
	}

	delete(ms.bookStorage, id)

	return nil

}
