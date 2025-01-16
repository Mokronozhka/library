package server

import (
	"github.com/gin-gonic/gin"
	"library/internal/domain/models"
	"library/internal/logger"
	"net/http"
)

func (s *ServerStruct) LoginUserHandler(ctx *gin.Context) {

	var user models.UserLoginStruct
	var UID, tokenString string
	var err error

	log := logger.Get()

	if err = ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Error().Err(err).Msg("loginHandler / Bad structure")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // ??? Почему тут мы пишем в ctx ответ?
		// Да ещё и с методом JSON?
		// В Вашем примере в gin.H передаётся err, а IDE мне предлагает err.Error(). В чём разница?
		return
	}

	if err = s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("loginHandler / Bad data")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // ??? Почему тут мы пишем в ctx ответ?
		return
	}

	UID, err = s.uService.LoginUser(user) // ??? Почему в LoginUser мы обращаемся через storage к ValidateUser(user)
	// а здесь сразу к LoginUser, пропуская storage? Хотя используем одинаковый обюъект UserServiceStruct.
	// Почему это возможно?

	if err != nil {
		log.Error().Err(err).Msg("loginHandler / Login fail")
		// Что за конструкция? Почему передаём Err и потом ещё добавляем MSG?
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid data", "error": err.Error()}) // Почему тут msg, а выше error?
		return
	}

	tokenString, err = CreateToken(UID)

	if err != nil {
		log.Error().Err(err).Msg("loginHandler / Token creation error")
		// Что за конструкция? Почему передаём Err и потом ещё добавляем MSG?
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Token creation error", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func (s *ServerStruct) RegistrationUserHandler(ctx *gin.Context) {

	var user models.UserStruct
	var UID, tokenString string
	var err error

	log := logger.Get()

	if err = ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Error().Err(err).Msg("RegistrationUserHandler / Bad structure")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("RegistrationUserHandler / Bad data")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UID, err = s.uService.RegistrationUser(user)

	if err != nil {
		log.Error().Err(err).Msg("RegistrationUserHandler / Login fail")
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid data", "error": err.Error()})
		return
	}

	tokenString, err = CreateToken(UID)

	if err != nil {
		log.Error().Err(err).Msg("RegistrationUserHandler / Token creation error")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Token creation error", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})

}
