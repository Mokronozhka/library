package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"library/internal/domain/models"
	"library/internal/logger"
	"library/internal/storage/storageerror"
	"time"
)

type DBStorage struct {
	conn *pgx.Conn
}

func NewDBStorage(ctx context.Context, addr string) (*DBStorage, error) {

	conn, err := pgx.Connect(ctx, addr)

	if err != nil {
		return nil, err
	}

	return &DBStorage{conn: conn}, nil

}

func (db *DBStorage) SaveUser(user models.UserStruct) (string, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost) // ??? Опять IDE сама всю строку сделала. Как?

	if err != nil {
		return "", err
	}

	user.Password = string(hash)

	user.ID = uuid.New() // Нужно генерировать в самой бд как-то

	_, err = db.conn.Exec(ctx, "INSERT INTO Users (ID, Name, Password, Email, Age) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Name, user.Password, user.Email, user.Age)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return "", storageerror.ErrUserAlreadyExist
			}

		}

		log.Error().Err(err).Msg("failed save user")

		return "", err

	}

	return user.ID.String(), nil

}

func (db *DBStorage) ValidateUser(user models.UserLoginStruct) (string, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	row := db.conn.QueryRow(ctx, "SELECT ID, Email, Password WHERE email = $1", user.Email)

	var userDB models.UserStruct

	if err := row.Scan(&userDB.ID, &userDB.Email, &userDB.Password); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return "", storageerror.ErrUserNotFound
		}

		log.Error().Err(err).Msg("failed validate user")
		return "", err

	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
		return "", storageerror.ErrUserInvalidPassword
	}

	return userDB.ID.String(), nil

}

func (db *DBStorage) GetBooks() ([]models.BookStruct, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	rows, err := db.conn.Query(ctx, "SELECT * FROM books")

	if err != nil {
		log.Error().Err(err).Msg("failed get data from table books")
	}

	var books []models.BookStruct

	for rows.Next() {

		var book models.BookStruct

		if err = rows.Scan(&book.ID, &book.Name, &book.Description, &book.Author, &book.DateWriting); err != nil {
			log.Error().Err(err).Msg("failed scan rows data")
			return nil, err
		}

		books = append(books, book)

	}

	return books, nil

}

func Migrations(dbDSN string, migratePath string) error {

	log := logger.Get()

	migrPath := fmt.Sprintf("file://%s", migratePath)

	m, err := migrate.New(migrPath, dbDSN)

	if err != nil {
		log.Error().Err(err).Msg("failed migrate one")
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Error().Err(err).Msg("failed migrate two")
			return nil
		}
		log.Error().Err(err).Msg("failed migrate three")
		return err
	}

	log.Debug().Msg("migrate ok")
	return nil

}
