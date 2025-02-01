package main

import (
	"context"
	"errors"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/sync/errgroup"
	"library/internal/config"
	"library/internal/logger"
	"library/internal/server"
	"library/internal/service"
	"library/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.ReadConfig()

	log := logger.Get(cfg.Debug)
	//log.Debug().Any("cfg", cfg).Msg("test cfg")
	//log.Info().Msg("info")
	//log.Warn().Msg("warn")
	//log.Error().Msg("error")
	//log.Fatal().Msg("fatal")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c // Блокируется выполнение, т.к. в канале нет ничего. Почитать про каналы ещё.
		log.Info().Msg("graceful shutdown")
		cancel()
	}()

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

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err = s.Run(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Send()
				return err
			}
		}
		return nil
	})

	group.Go(func() error {
		log.Debug().Msg("start listening error channel")
		defer log.Debug().Msg("stop listening error channel")
		return <-s.ChanErr
	})

	group.Go(func() error {
		<-gCtx.Done()
		return s.Shutdown(gCtx)
	})

	group.Go(func() error {
		<-gCtx.Done()
		return DBStorage.Close()
	})

	if err = group.Wait(); err != nil {
		log.Error().Err(err).Send()
	}

	log.Info().Msg("shutdown complete")

}
