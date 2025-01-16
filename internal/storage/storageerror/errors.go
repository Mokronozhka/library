package storageerror

import "errors"

var (
	ErrBookAlreadyExist = errors.New("book already exists")
	ErrBookStorageEmpty = errors.New("book storage is empty")
	ErrBookNotFound     = errors.New("book not found")
)
