package service

import (
	"bytes"
	"encoding/json"
	"library/internal/domain/models"
	"library/internal/logger"
)

type UserStorage interface {
	SaveUser(models.UserStruct) (string, error)
	ValidateUser(models.UserLoginStruct) (string, error)
}

type UserServiceStruct struct {
	storage UserStorage
}

func NewUserService(storage UserStorage) UserServiceStruct {

	return UserServiceStruct{storage: storage}

}

func (us UserServiceStruct) LoginUser(user models.UserLoginStruct) (string, error) {

	log := logger.Get()

	json.NewDecoder(bytes.NewBufferString(user.Password)).Decode(&user)

	UID, err := us.storage.ValidateUser(user)

	if err != nil {
		log.Error().Err(err).Msg("LoginUser / User validation error")
		return "", err
	}

	return UID, nil

}

func (us UserServiceStruct) RegistrationUser(user models.UserStruct) (string, error) {

	log := logger.Get()

	UID, err := us.storage.SaveUser(user)

	if err != nil {
		log.Error().Err(err).Msg("RegistrationUser / User save error")
		return "", err
	}

	return UID, nil

}
