package storage

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Storage struct {
	config           *Config
	db               *sql.DB
	userRepository   *UserRepository
	dialogRepository *DialogRepository
	jokeRepository   *JokeRepository
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (storage *Storage) Open() error {
	dataSourceName := storage.config.User + ":" + storage.config.Password + "@/" + storage.config.DataBase
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	storage.db = db
	log.Println("Database connection has been set up successfully")
	return nil
}

func (storage *Storage) Close() {
	if err := storage.db.Close(); err != nil {
		log.Println("Error during closing DB connection: ", err)
	}
}

func (storage *Storage) User() *UserRepository {
	if storage.userRepository != nil {
		return storage.userRepository
	}
	storage.userRepository = &UserRepository{
		storage: storage,
	}
	return storage.userRepository
}

func (storage *Storage) Dialog() *DialogRepository {
	if storage.dialogRepository != nil {
		return storage.dialogRepository
	}
	storage.dialogRepository = &DialogRepository{
		storage: storage,
	}
	return storage.dialogRepository
}

func (storage *Storage) Joke() *JokeRepository {
	if storage.jokeRepository != nil {
		return storage.jokeRepository
	}
	storage.jokeRepository = &JokeRepository{
		storage: storage,
	}
	return storage.jokeRepository
}
