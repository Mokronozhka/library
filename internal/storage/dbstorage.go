package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
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

func Migrations(dbDSN string, migratePath string) error {

	log := logger.Get()

	migrationPath := fmt.Sprintf("file://%s", migratePath)

	var err error
	var m *migrate.Migrate

	m, err = migrate.New(migrationPath, dbDSN)

	if err != nil {
		log.Error().Err(err).Msg("failed migrate one")
		return err
	}

	if err = m.Up(); err != nil {

		if !errors.Is(err, migrate.ErrNoChange) {
			log.Error().Err(err).Msg("failed migrate three")
			return nil
		}

	}

	log.Debug().Msg("migrate ok")
	return nil

}

func (db *DBStorage) Close() error {
	return db.conn.Close(context.Background())

}

func (db *DBStorage) GetUsers() ([]models.UserStruct, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	rows, err := db.conn.Query(ctx, "SELECT ID, Name, Password, Email, Age, DateRegistration FROM Users")

	if err != nil {
		log.Error().Err(err).Msg("Failed get data from table Users")
		return nil, err
	}

	var users []models.UserStruct

	for rows.Next() {

		var user models.UserStruct

		if err = rows.Scan(&user.ID,
			&user.Name,
			&user.Password,
			&user.Email,
			&user.Age,
			&user.DateRegistration); err != nil {
			log.Error().Err(err).Msg("Failed scan rows data")
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil

}

func (db *DBStorage) GetUser(id string) (models.UserStruct, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	userDB := models.UserStruct{}

	ID, err := uuid.Parse(id)

	if err != nil {
		log.Error().Err(err).Msg("Failed parse ID")
		return userDB, err
	}

	row := db.conn.QueryRow(ctx,
		"SELECT ID, Name, Password, Email, Age, DateRegistration FROM Users WHERE ID = $1", ID)

	if err = row.Scan(&userDB.ID,
		&userDB.Name,
		&userDB.Password,
		&userDB.Email,
		&userDB.Age,
		&userDB.DateRegistration); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return userDB, storageerror.ErrUserNotFound
		}

		log.Error().Err(err).Msg("Failed get data from table Users")
		return userDB, err

	}

	return userDB, nil

}

func (db *DBStorage) SaveUser(user models.UserStruct) (string, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	row := db.conn.QueryRow(ctx,
		"SELECT ID FROM Users WHERE Email = $1", user.Email)

	var IDTemp uuid.UUID

	if err := row.Scan(&IDTemp); err != nil {

		if !errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("Failed get data from table Books")
			return "", err
		}

	} else {
		return "", storageerror.ErrUserAlreadyExist
	}

	user.ID = uuid.New() // Нужно генерировать в самой бд как-то

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	user.Password = string(hash)

	_, err = db.conn.Exec(ctx, "INSERT INTO Users (ID, Name, Password, Email, Age) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Name, user.Password, user.Email, user.Age)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return "", storageerror.ErrUserAlreadyExist
			}

		}

		log.Error().Err(err).Msg("Failed save user")

		return "", err

	}

	return user.ID.String(), nil

}

func (db *DBStorage) ValidateUser(user models.UserLoginStruct) (string, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	row := db.conn.QueryRow(ctx, "SELECT ID, Password FROM Users WHERE email = $1", user.Email)

	var userDB models.UserStruct

	if err := row.Scan(&userDB.ID, &userDB.Password); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return "", storageerror.ErrUserNotFound
		}

		log.Error().Err(err).Msg("Failed validate user")
		return "", err

	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
		return "", storageerror.ErrUserInvalidPassword
	}

	return userDB.ID.String(), nil

}

func (db *DBStorage) EditUser(id string, user models.UserStruct) error {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	ID, err := uuid.Parse(id)

	if err != nil {
		log.Error().Err(err).Msg("Failed parse ID")
		return err
	}

	var userDB models.UserStruct

	row := db.conn.QueryRow(ctx,
		"SELECT ID, Password, Email FROM Users WHERE ID = $1", ID)

	if err = row.Scan(&userDB.ID,
		&userDB.Password,
		&userDB.Email); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return storageerror.ErrUserNotFound
		}

		log.Error().Err(err).Msg("Failed get data from table Users")
		return err

	}

	if userDB.Email != user.Email {

		row = db.conn.QueryRow(ctx,
			"SELECT ID FROM Users WHERE Email = $1", user.Email)

		var IDTemp uuid.UUID

		if err = row.Scan(&IDTemp); err != nil {

			if !errors.Is(err, pgx.ErrNoRows) {
				log.Error().Err(err).Msg("Failed get data from table Users")
				return err
			}

		} else {
			return fmt.Errorf("user with email '%s' already exists", user.Email)
		}

	}

	if errCompare := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); errCompare != nil {

		hash, errGen := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

		if errGen != nil {
			return errGen
		}

		user.Password = string(hash)

	} else {
		user.Password = userDB.Password
	}

	_, err = db.conn.Exec(ctx, "UPDATE Users SET Name = $1, Password = $2, Email = $3, Age = $4 WHERE ID = $5",
		user.Name, user.Password, user.Email, user.Age, userDB.ID)

	if err != nil {
		log.Error().Err(err).Msg("Failed edit user")
		return err
	}

	return nil

}

func (db *DBStorage) DeleteUser(id string) error {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	ID, err := uuid.Parse(id)

	if err != nil {
		log.Error().Err(err).Msg("Failed parse ID")
		return err
	}

	row := db.conn.QueryRow(ctx, "SELECT ID FROM Users WHERE ID = $1", ID)

	var IDTemp uuid.UUID

	if err = row.Scan(&IDTemp); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return storageerror.ErrUserNotFound
		}

		log.Error().Err(err).Msg("Failed get data from table Users")
		return err

	}

	_, err = db.conn.Exec(ctx, "DELETE FROM Users WHERE ID = $1", ID)

	if err != nil {
		log.Error().Err(err).Msg("Failed delete user")
		return err
	}

	return nil

}

func (db *DBStorage) GetBooks() ([]models.BookStruct, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	rows, err := db.conn.Query(ctx, "SELECT ID, Name, Description, Author, DateWriting FROM Books")

	if err != nil {
		log.Error().Err(err).Msg("failed get data from table Books")
		return nil, err
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

func (db *DBStorage) GetBook(id string) (models.BookStruct, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	bookDB := models.BookStruct{}

	ID, err := uuid.Parse(id)

	if err != nil {
		log.Error().Err(err).Msg("Failed parse ID")
		return bookDB, err
	}

	row := db.conn.QueryRow(ctx,
		"SELECT ID, Name, Description, Author, DateWriting FROM Books WHERE ID = $1", ID)

	if err = row.Scan(&bookDB.ID,
		&bookDB.Name,
		&bookDB.Description,
		&bookDB.Author,
		&bookDB.DateWriting); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return bookDB, storageerror.ErrBookNotFound
		}

		log.Error().Err(err).Msg("Failed get data from table Books")
		return bookDB, err

	}

	return bookDB, nil

}

func (db *DBStorage) SaveBook(book models.BookStruct) (string, error) {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	row := db.conn.QueryRow(ctx,
		"SELECT ID FROM Books WHERE Name = $1 and Author = $2", book.Name, book.Author)

	var IDTemp uuid.UUID

	if err := row.Scan(&IDTemp); err != nil {

		if !errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("Failed get data from table Books")
			return "", err
		}

	} else {
		return "", storageerror.ErrBookAlreadyExist
	}

	book.ID = uuid.New()

	_, err := db.conn.Exec(ctx,
		"INSERT INTO Books (ID, Name, Description, Author, DateWriting) VALUES ($1, $2, $3, $4, $5)",
		book.ID, book.Name, book.Description, book.Author, book.DateWriting)

	if err != nil {
		log.Error().Err(err).Msg("Failed save book")
		return "", err
	}

	return book.ID.String(), nil

}

func (db *DBStorage) EditBook(id string, book models.BookStruct) error {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	ID, err := uuid.Parse(id)

	if err != nil {
		log.Error().Err(err).Msg("Failed parse ID")
		return err
	}

	var bookDB models.BookStruct

	row := db.conn.QueryRow(ctx,
		"SELECT ID, Name, Author FROM Books WHERE ID = $1", ID)

	if err = row.Scan(&bookDB.ID,
		&bookDB.Name,
		&bookDB.Author); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return storageerror.ErrBookNotFound
		}

		log.Error().Err(err).Msg("Failed get data from table Books")
		return err

	}

	if bookDB.Name != book.Name || bookDB.Author != book.Author {

		row = db.conn.QueryRow(ctx,
			"SELECT ID FROM Books WHERE Name = $1 AND Author=$2", book.Name, bookDB.Author)

		var IDTemp uuid.UUID

		if err = row.Scan(&IDTemp); err != nil {

			if !errors.Is(err, pgx.ErrNoRows) {
				log.Error().Err(err).Msg("Failed get data from table Books")
				return err
			}

		} else {
			return fmt.Errorf("book with name '%s' and author '%s' already exists", book.Name, bookDB.Author)
		}

	}

	_, err = db.conn.Exec(ctx,
		"UPDATE Books SET Name = $1, Description = $2, Author = $3,  DateWriting =$4 WHERE ID = $5",
		book.Name, book.Description, book.Author, book.DateWriting, bookDB.ID)

	if err != nil {
		log.Error().Err(err).Msg("Failed edit book")
		return err
	}

	return nil

}

func (db *DBStorage) DeleteBook(id string) error {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	ID, err := uuid.Parse(id)

	if err != nil {
		log.Error().Err(err).Msg("Failed parse ID")
		return err
	}

	row := db.conn.QueryRow(ctx, "SELECT ID FROM Books WHERE ID = $1", ID)

	var IDTemp uuid.UUID

	if err = row.Scan(&IDTemp); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return storageerror.ErrBookNotFound
		}

		log.Error().Err(err).Msg("Failed get data from table Books")
		return err

	}

	_, err = db.conn.Exec(ctx, "UPDATE Books SET Deleted = true WHERE ID = $1", ID)

	if err != nil {
		log.Error().Err(err).Msg("Failed delete book")
		return err
	}

	return nil

}

func (db *DBStorage) DeleteBooks() error {

	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	tx, err := db.conn.Begin(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Failed create transaction")
		return err
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed rollback transaction")
		}
	}()

	_, err = db.conn.Exec(ctx, "DELETE FROM Books WHERE Deleted = true")

	if err != nil {
		log.Error().Err(err).Msg("Failed delete book")
		return err
	}

	return tx.Commit(ctx)

}
