package main

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"library/internal/config"
	"library/internal/logger"
	"library/internal/server"
	"library/internal/service"
	"library/internal/storage"
)

func main() {

	cfg := config.ReadConfig()

	log := logger.Get(cfg.Debug)
	//log.Debug().Any("cfg", cfg).Msg("test cfg")
	//log.Info().Msg("info")
	//log.Warn().Msg("warn")
	//log.Error().Msg("error")
	//log.Fatal().Msg("fatal")

	var err error

	err = storage.Migrations(cfg.DbDSN, cfg.MigratePath)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	//storage := storage.NewMapStorage()

	//userService := service.NewUserService(storage) // Непонятно. Мы сюда передаём структуру из мапов,
	//// а внутри идёт параметр с типом UserStorage. Как так? И что мы в итоге возвращаем?
	//bookService := service.NewBookService(storage)

	var DBStorage *storage.DBStorage
	var MapStorage *storage.MapStorage
	var userService service.UserServiceStruct
	var bookService service.BookServiceStruct

	DBStorage, err = storage.NewDBStorage(context.Background(), cfg.DbDSN)

	if err != nil {

		log.Error().Err(err).Send()

		MapStorage = storage.NewMapStorage()
		userService = service.NewUserService(MapStorage)
		bookService = service.NewBookService(MapStorage)

	} else {

		userService = service.NewUserService(DBStorage)
		bookService = service.NewBookService(DBStorage)

	}

	s := server.New(cfg, userService, bookService)

	if err := s.Run(); err != nil {
		log.Fatal().Err(err).Send()
	}

}
