package models

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	TelegramId *int   `json:"telegram_id"`
	Greeting   struct {
		ID      int    `json:"id"`
		UserId  int    `json:"user_id"`
		Message string `json:"message"`
	} `json:"greeting"`
}
