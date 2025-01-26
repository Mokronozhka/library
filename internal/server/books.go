package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"library/internal/domain/models"
	"library/internal/logger"
	"library/internal/storage/storageerror"
	"net/http"
	"time"
)

func (s *ServerStruct) GetBooksHandler(ctx *gin.Context) {

	log := logger.Get()

	books, err := s.bService.GetBooks()

	if err != nil {
		log.Error().Err(err).Msg("Get books failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": books})

}

func (s *ServerStruct) GetBookHandler(ctx *gin.Context) {

	log := logger.Get()

	id := ctx.Param("id")

	if id == "" {
		log.Error().Msg("Book ID is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Book ID is empty"})
		return
	}

	book, err := s.bService.GetBook(id)

	if err != nil {
		log.Error().Err(err).Msg("Get book failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": book})

}

func (s *ServerStruct) AddBookHandler(ctx *gin.Context) {

	log := logger.Get()

	var book models.BookStruct

	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		log.Error().Err(err).Msg("Unmarshall body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.valid.Struct(book); err != nil {
		log.Error().Err(err).Msg("Invalid body structure")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.bService.AddBook(book)

	if err != nil {

		log.Error().Err(err).Msg("Save book failed")

		status := http.StatusInternalServerError

		if errors.Is(err, storageerror.ErrBookAlreadyExist) {
			status = http.StatusConflict
		}

		ctx.JSON(status, gin.H{"error": err.Error()})

		return

	}

	ctx.JSON(http.StatusCreated, gin.H{"result": fmt.Sprintf("Book added. ID - %s", id)})

}

func (s *ServerStruct) EditBookHandler(ctx *gin.Context) {

	log := logger.Get()

	id := ctx.Param("id")

	if id == "" {
		log.Error().Msg("Book ID is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Book ID is empty"})
		return
	}

	var book models.BookStruct

	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		log.Error().Err(err).Msg("Unmarshall body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.valid.Struct(book); err != nil {
		log.Error().Err(err).Msg("Invalid body structure")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.bService.EditBook(id, book)

	if err != nil {
		log.Error().Err(err).Msg("Edit book failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "Book edited"})

}

func (s *ServerStruct) DeleteBookHandler(ctx *gin.Context) {

	log := logger.Get()

	id := ctx.Param("id")

	if id == "" {
		log.Error().Msg("Book ID is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Book ID is empty"})
		return
	}

	err := s.bService.DeleteBook(id)

	if err != nil {
		log.Error().Err(err).Msg("Delete book failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.chanDel <- struct{}{}
	ctx.JSON(http.StatusOK, gin.H{"result": "Book removed"})

}

func (s *ServerStruct) deleter(ctx context.Context) {
	log := logger.Get()
	defer log.Debug().Msg("deleter stopped")
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("deleter. ctx.Done()")
			return
		case <-time.After(5 * time.Minute):
			//case <-time.After(5 * time.Second):
			//log.Debug().Int("len", len(s.chanDel)).Int("cap", cap(s.chanDel)).Msg("deleter. channel check")
			if len(s.chanDel) == cap(s.chanDel) {
				log.Debug().Int("len", len(s.chanDel)).Int("cap", cap(s.chanDel)).Msg("deleter. channel full")
				for i := 0; i < cap(s.chanDel); i++ {
					<-s.chanDel
				}
				if err := s.bService.DeleteBooks(); err != nil {
					log.Error().Err(err).Msg("Delete books failed")
					s.ChanErr <- err
					return
				}
			}

		}
	}
}
