package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"library/internal/config"
	"library/internal/logger"
	"library/internal/service"
	"net/http"
)

//type Storage interface {
//	SaveUser()
//	ValidUser()
//	SaveBook()
//	DeleteBook()
//	GetBooks()
//	GetBook()
//}

type ServerStruct struct {
	server   *http.Server              // Структура Server из пакета http, чтобы мы могли вызывать грейсфул шатдаун
	valid    *validator.Validate       // Ссылка на оригинальную переменную
	uService service.UserServiceStruct // Копия
	bService service.BookServiceStruct // Копия
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
	}

	// ??? Почему мы валидатор и структуры пользователя и книги запихиваем в структуру сервера?

}

func (s *ServerStruct) Run() error {

	log := logger.Get()

	router := s.configRouting() // Конфигурацию путей вынесли в отдельную функцию
	s.server.Handler = router

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
		users.POST("/login", s.LoginUserHandler) // ??? В параметрах функции LoginUserHandler у нас контекст, но здесь мы его не указываем
		// Как он тогда передаётся туда? И почему мы указываем функцию без ()?
		//users.POST("/logout")
		users.GET("/info")
	}

	books := router.Group("/books")
	{
		books.GET("/", s.GetBooksHandler)
		books.GET("/:id", s.GetBookHandler)
		books.POST("/add", s.AddBookHandler)
		//books.POST("/edit")
		books.DELETE("/:id", s.DeleteBookHandler)
	}

	return router

}
