package storage

import (
	"database/sql"
	"fmt"
	"github.com/drTragger/rileyBot/internal/app/models"
)

type UserRepository struct {
	storage *Storage
}

var (
	tableUsers     = "users"
	tableGreetings = "greetings"
)

func (ur *UserRepository) Create(u *models.User) error {
	query := fmt.Sprintf("INSERT INTO %s (username, telegram_id) VALUES (?, ?)", tableUsers)
	if _, err := ur.storage.db.Query(query, u.Username, u.TelegramId); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) UserExists(username string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE username=?)", tableUsers)
	row := ur.storage.db.QueryRow(query, username)
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (ur *UserRepository) FindByTelegramUsernameWithGreetings(username string) (*models.User, bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s RIGHT JOIN %s ON %s.id = %s.user_id WHERE %s.username=?", tableUsers, tableGreetings, tableUsers, tableGreetings, tableUsers)
	user := models.User{}
	row := ur.storage.db.QueryRow(query, username)
	switch err := row.Scan(&user.ID, &user.TelegramId, &user.Username, &user.Greeting.ID, &user.Greeting.UserId, &user.Greeting.Message); err {
	case sql.ErrNoRows:
		return nil, false, nil
	case nil:
		return &user, true, nil
	default:
		return nil, false, err
	}
}

func (ur *UserRepository) FindByTelegramUsername(username string) (*models.User, bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=?", tableUsers)
	user := models.User{}
	row := ur.storage.db.QueryRow(query, username)
	switch err := row.Scan(&user.ID, &user.TelegramId, &user.Username); err {
	case sql.ErrNoRows:
		return nil, false, nil
	case nil:
		return &user, true, nil
	default:
		return nil, false, err
	}
}

func (ur *UserRepository) FindById(id int) (*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", tableUsers)
	user := models.User{}
	row := ur.storage.db.QueryRow(query, id)
	if err := row.Scan(&user.ID, &user.TelegramId, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}
