package riley

import (
	"github.com/drTragger/RileyBot/internal/app/models"
	"github.com/yanzay/tbot/v2"
	"log"
	"strconv"
)

const errorMessage = "Вибачте, тимчасово туплю.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Ножиці-Папір /play"

func handleChatActionError(err error) {
	if err != nil {
		log.Println("Error chat action: ", err.Error())
	}
}

func handleMessageError(message *tbot.Message, err error) {
	if err != nil {
		log.Printf("Message: %s\nError: %s", message.Text, err.Error())
	}
}

func (b *Bot) NewDialog(username string, chatId string, dialogName string) (string, bool) {
	var user *models.User
	user, ok, err := b.storage.User().FindByTelegramUsername(username)

	if err != nil {
		b.logger.Info("Error during fetching user data: ", err.Error())
		return errorMessage, false
	}
	if !ok {
		userId, err := strconv.Atoi(chatId)
		if err != nil {
			b.logger.Info("Failed to convert user ID ", err.Error())
			return errorMessage, false
		}

		user = &models.User{Username: username, TelegramId: &userId}
		err = b.storage.User().Create(user)
		if err != nil {
			b.logger.Info("Failed to create new user: ", err.Error())
			return errorMessage, false
		}
	}
	err = b.storage.Dialog().Create(&models.Dialog{Name: dialogName, UserId: user.ID, Status: true})
	if err != nil {
		b.logger.Info("Failed to create new dialog: ", err.Error())
		return errorMessage, false
	}
	return "", true
}

func (b *Bot) CheckDialogStatus(username string, dialogName string) (string, *models.Dialog, bool) {
	user, ok, err := b.storage.User().FindByTelegramUsername(username)

	if err != nil {
		b.logger.Info("Error during fetching user data: ", err.Error())
		return errorMessage, nil, false
	}

	if !ok {
		b.logger.Info("User and dialog not found")
		return "Будь ласка, запустіть мене, виконавши команду /start", nil, false
	}

	dialog, ok, err := b.storage.Dialog().FindLatestUserDialog(user.ID, dialogName)
	if err != nil {
		b.logger.Error("Error during fetching dialog data: ", err.Error())
		return errorMessage, dialog, false
	}

	if !ok || dialog.Status != true {
		b.logger.Info("No active dialog status")
		return "Перепрошую, я поки не вмію розпізнавати такі повідомлення. Спробуйте:\n\n/play - Пограти у Камінь-Ножиці-Папір\n\n/weather - Дізнатись, яка зараз погода", dialog, false
	}
	return "", dialog, true
}
