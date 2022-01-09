package models

type Dialog struct {
	ID     int    `json:"id"`
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
}
