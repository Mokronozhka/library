package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"library/internal/domain/models"
	"library/internal/logger"
	"library/internal/storage/storageerror"
	"net/http"
)

func (s *ServerStruct) GetBooksHandler(ctx *gin.Context) {

	log := logger.Get()

	books, err := s.bService.GetBooks()

	if err != nil {

		log.Error().Err(err).Msg("Get books failed")

		if errors.Is(err, storageerror.ErrBookStorageEmpty) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}

	ctx.JSON(http.StatusOK, gin.H{"result": books})

}

func (s *ServerStruct) GetBookHandler(ctx *gin.Context) {

	log := logger.Get()

	id := ctx.Param("id")

	if id == "" {
		log.Error().Msg("Book id is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "book id is empty"})
		return
	}

	book, err := s.bService.GetBook(id)

	if err != nil {

		log.Error().Err(err).Msg("Get book failed")

		if errors.Is(err, storageerror.ErrBookNotFound) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}

	ctx.JSON(http.StatusCreated, gin.H{"result": book})

}

func (s *ServerStruct) AddBookHandler(ctx *gin.Context) {

	log := logger.Get()

	var book models.BookStruct

	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		log.Error().Err(err).Msg("Unmarshall body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := s.bService.AddBook(book)

	if err != nil {

		log.Error().Err(err).Msg("Save book failed")

		if errors.Is(err, storageerror.ErrBookAlreadyExist) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}

	ctx.JSON(http.StatusCreated, gin.H{"result": "book added"})

}

func (s *ServerStruct) DeleteBookHandler(ctx *gin.Context) {

	log := logger.Get()

	id := ctx.Param("id")

	if id == "" {
		log.Error().Msg("Book id is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "book id is empty"})
		return
	}

	err := s.bService.DeleteBook(id)

	if err != nil {

		log.Error().Err(err).Msg("Get book failed")

		if errors.Is(err, storageerror.ErrBookNotFound) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}

	ctx.JSON(http.StatusOK, gin.H{"result": "book removed"})

}
