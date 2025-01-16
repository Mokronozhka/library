package main

import (
	"library/internal/config"
	"library/internal/logger"
	"library/internal/server"
	"library/internal/service"
	"library/internal/storage"
)

func main() {

	cfg := config.ReadConfig()
	//fmt. (cfg)

	log := logger.Get(cfg.Debug) // Почему каждый раз вызываем гет?
	//log.Debug().Any("cfg", cfg).Msg("test cfg")
	//log.Info().Msg("info")
	//log.Warn().Msg("warn")
	//log.Error().Msg("error")
	//log.Fatal().Msg("fatal")

	mapStorage := storage.NewMapStorage()

	userService := service.NewUserService(mapStorage) // Непонятно. Мы сюда передаём структуру из мапов,
	// а внутри идёт параметр с типом UserStorage. Как так? И что мы в итоге возвращаем?
	bookService := service.NewBookService(mapStorage)

	s := server.New(cfg, userService, bookService)

	if err := s.Run(); err != nil { // Чисто запуск сервера
		log.Fatal().Err(err).Send() // ??? Почему тут используем Send, а не MSG, как в server.Run()?
	}

}
