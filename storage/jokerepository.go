package storage

import (
	"database/sql"
	"fmt"
	"github.com/drTragger/RileyBot/internal/app/models"
)

type JokeRepository struct {
	storage *Storage
}

var (
	tableJokes = "jokes"
)

func (jr *JokeRepository) GetRandom() (*models.Joke, bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY RAND() LIMIT 1", tableJokes)
	joke := models.Joke{}
	row := jr.storage.db.QueryRow(query)
	switch err := row.Scan(&joke.ID, &joke.Joke); err {
	case sql.ErrNoRows:
		return nil, false, nil
	case nil:
		return &joke, true, nil
	default:
		return nil, false, err
	}
}
