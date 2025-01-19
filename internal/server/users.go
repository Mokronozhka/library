package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"library/internal/domain/models"
	"library/internal/logger"
	"library/internal/server/utils"
	"library/internal/storage/storageerror"
	"net/http"
)

func (s *ServerStruct) RegistrationUserHandler(ctx *gin.Context) {

	log := logger.Get()

	var user models.UserStruct
	var err error
	var ID, token string

	if err = ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Error().Err(err).Msg("Unmarshall body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("Invalid body structure")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID, err = s.uService.RegistrationUser(user)

	if err != nil {
		log.Error().Err(err).Msg("Registration user fail")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err = util.CreateToken(ID)

	if err != nil {
		log.Error().Err(err).Msg("Token creation error")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", token)

	ctx.JSON(http.StatusOK,
		gin.H{"result": fmt.Sprintf("User registered. ID - %s", ID)})

}

func (s *ServerStruct) LoginUserHandler(ctx *gin.Context) {

	log := logger.Get()

	var user models.UserLoginStruct
	var err error
	var ID, token string

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Error().Err(err).Msg("Unmarshall body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("Invalid body structure")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID, err = s.uService.LoginUser(user)

	if err != nil {
		log.Error().Err(err).Msg("Login user fail")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err = util.CreateToken(ID)

	if err != nil {
		log.Error().Err(err).Msg("Token creation error")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", token)

	ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprintf("User logged. ID - %s", ID)})

}

func (s *ServerStruct) GetUsersHandler(ctx *gin.Context) {

	log := logger.Get()

	users, err := s.uService.GetUsers()

	if err != nil {
		log.Error().Err(err).Msg("get users error")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": users})

}

func (s *ServerStruct) GetUserHandler(ctx *gin.Context) {

	log := logger.Get()

	id := ctx.Param("id")

	if id == "" {
		log.Error().Msg("User ID is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is empty"})
		return
	}

	user, err := s.uService.GetUser(id)

	if err != nil {
		log.Error().Err(err).Msg("Get user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": user})

}

func (s *ServerStruct) AddUserHandler(ctx *gin.Context) {

	log := logger.Get()

	var user models.UserStruct

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Error().Err(err).Msg("Unmarshall body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("Invalid body structure")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.uService.AddUser(user)

	if err != nil {

		log.Error().Err(err).Msg("Add user failed")

		status := http.StatusInternalServerError

		if errors.Is(err, storageerror.ErrUserAlreadyExist) {
			status = http.StatusConflict
		}

		ctx.JSON(status, gin.H{"error": err.Error()})

		return

	}

	ctx.JSON(http.StatusCreated, gin.H{"result": fmt.Sprintf("User added. ID - %s", id)})

}

func (s *ServerStruct) EditUserHandler(ctx *gin.Context) {

	log := logger.Get()

	id := ctx.Param("id")

	if id == "" {
		log.Error().Msg("User ID is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is empty"})
		return
	}

	var user models.UserStruct

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Error().Err(err).Msg("Unmarshall body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("Invalid body structure")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.uService.EditUser(id, user)

	if err != nil {
		log.Error().Err(err).Msg("Edit user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "User edited"})

}

func (s *ServerStruct) DeleteUserHandler(ctx *gin.Context) {

	log := logger.Get()

	id := ctx.Param("id")

	if id == "" {
		log.Error().Msg("User ID is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is empty"})
		return
	}

	err := s.uService.DeleteUser(id)

	if err != nil {
		log.Error().Err(err).Msg("Delete user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "User removed"})

}
