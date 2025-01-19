package service

import "library/internal/domain/models"

type UserStorage interface {
	GetUsers() ([]models.UserStruct, error)
	GetUser(string) (models.UserStruct, error)
	SaveUser(models.UserStruct) (string, error)
	ValidateUser(models.UserLoginStruct) (string, error)
	EditUser(string, models.UserStruct) error
	DeleteUser(string) error
}

type UserServiceStruct struct {
	storage UserStorage
}

func NewUserService(storage UserStorage) UserServiceStruct {
	return UserServiceStruct{storage: storage}
}

func (us UserServiceStruct) RegistrationUser(user models.UserStruct) (string, error) {
	return us.storage.SaveUser(user)
}

func (us UserServiceStruct) LoginUser(user models.UserLoginStruct) (string, error) {
	return us.storage.ValidateUser(user)
}

func (us UserServiceStruct) GetUsers() ([]models.UserStruct, error) {
	return us.storage.GetUsers()
}

func (us UserServiceStruct) GetUser(id string) (models.UserStruct, error) {
	return us.storage.GetUser(id)
}

func (us UserServiceStruct) AddUser(user models.UserStruct) (string, error) {
	return us.storage.SaveUser(user)
}

func (us UserServiceStruct) EditUser(id string, user models.UserStruct) error {
	return us.storage.EditUser(id, user)
}

func (us UserServiceStruct) DeleteUser(id string) error {
	return us.storage.DeleteUser(id)
}
