package storage

import (
	"database/sql"
	"fmt"
	"github.com/drTragger/RileyBot/internal/app/models"
)

type DialogRepository struct {
	storage *Storage
}

var (
	tableDialogs = "dialogs"
)

func (dr *DialogRepository) Create(d *models.Dialog) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, name, status) VALUES (?, ?, ?)", tableDialogs)
	if _, err := dr.storage.db.Query(query, d.UserId, d.Name, d.Status); err != nil {
		return err
	}
	return nil
}

func (dr *DialogRepository) UpdateStatus(id int) error {
	query := fmt.Sprintf("UPDATE %s SET status=0 WHERE id=?", tableDialogs)
	if _, err := dr.storage.db.Query(query, id); err != nil {
		return err
	}
	return nil
}

func (dr *DialogRepository) FindLatestUserDialog(userId int, dialogName string) (*models.Dialog, bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? AND name=? ORDER BY id DESC LIMIT 1", tableDialogs)
	dialog := models.Dialog{}
	row := dr.storage.db.QueryRow(query, userId, dialogName)
	switch err := row.Scan(&dialog.ID, &dialog.UserId, &dialog.Name, &dialog.Status); err {
	case sql.ErrNoRows:
		return nil, false, nil
	case nil:
		return &dialog, true, nil
	default:
		return nil, false, err
	}
}
