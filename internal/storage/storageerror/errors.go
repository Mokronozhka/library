package storageerror

import "errors"

var (
	ErrBookAlreadyExist = errors.New("book already exists")
	ErrBookStorageEmpty = errors.New("book storage is empty")
	ErrBookNotFound     = errors.New("book not found")

	ErrUserAlreadyExist    = errors.New("user already exists")
	ErrUserStorageEmpty    = errors.New("user storage is empty")
	ErrUserInvalidPassword = errors.New("user invalid password")
	ErrUserNotFound        = errors.New("user not found")
)
