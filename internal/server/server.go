package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"library/internal/config"
	"library/internal/logger"
	"library/internal/server/utils"
	"library/internal/service"
	"net/http"
)

type ServerStruct struct {
	server   *http.Server              // Структура Server из пакета http, чтобы мы могли вызывать грейсфул шатдаун
	valid    *validator.Validate       // Ссылка на оригинальную переменную
	uService service.UserServiceStruct // Копия
	bService service.BookServiceStruct // Копия
	chanDel  chan struct{}
	ChanErr  chan error
}

func New(cfg config.ConfigStruct,
	uService service.UserServiceStruct,
	bService service.BookServiceStruct) *ServerStruct {

	addrStr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := http.Server{
		Addr: addrStr,
	}

	valid := validator.New()

	return &ServerStruct{
		server:   &server, // ??? Почему тут &server?
		valid:    valid,   // ??? Почему тут без &?
		uService: uService,
		bService: bService,
		chanDel:  make(chan struct{}, 10),
		ChanErr:  make(chan error, 10),
	}

	// ??? Почему мы валидатор и структуры пользователя и книги запихиваем в структуру сервера?

}

func (s *ServerStruct) Run(ctx context.Context) error {

	log := logger.Get()

	router := s.configRouting() // Конфигурацию путей вынесли в отдельную функцию
	s.server.Handler = router

	go s.deleter(ctx)

	log.Info().Str("addr", s.server.Addr).Msg("starting server")

	if err := s.server.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("running server failed")
		return err
	}

	return nil

}

func (s *ServerStruct) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *ServerStruct) configRouting() *gin.Engine {

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})

	users := router.Group("/users")
	{
		users.POST("/registration", s.RegistrationUserHandler)
		users.POST("/login", s.LoginUserHandler)
		users.GET("/", s.JWTAuthMiddleware(), s.GetUsersHandler)
		users.GET("/:id", s.JWTAuthMiddleware(), s.GetUserHandler)
		users.POST("/", s.JWTAuthMiddleware(), s.AddUserHandler)
		users.PUT("/:id", s.JWTAuthMiddleware(), s.EditUserHandler)
		users.DELETE("/:id", s.JWTAuthMiddleware(), s.DeleteUserHandler)
	}

	books := router.Group("/books")
	{
		books.GET("/", s.JWTAuthMiddleware(), s.GetBooksHandler)
		books.GET("/:id", s.JWTAuthMiddleware(), s.GetBookHandler)
		books.POST("/", s.JWTAuthMiddleware(), s.AddBookHandler)
		books.PUT("/:id", s.JWTAuthMiddleware(), s.EditBookHandler)
		books.DELETE("/:id", s.JWTAuthMiddleware(), s.DeleteBookHandler)
	}

	return router

}

func (s *ServerStruct) JWTAuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		log := logger.Get()
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			ctx.Abort()
			return
		}
		err := util.ValidateToken(token)
		if err != nil {
			log.Error().Err(err).Send()
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Next() // Работает и без него
	}

}
