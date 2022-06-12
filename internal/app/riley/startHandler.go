package riley

import (
	"github.com/drTragger/RileyBot/internal/app/models"
	"github.com/yanzay/tbot/v2"
	"strconv"
	"time"
)

func (b *Bot) StartHandler(m *tbot.Message) {
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)
	var stdMessage = "Вітаю, я бот Райлі🖖\n\nСлава Україні🇺🇦\n\n/play\tКамінь-Ножиці-Папір\n\n/weather\tДізнатись прогноз погоди"
	var msg string
	userId, err := strconv.Atoi(m.Chat.ID)
	if err != nil {
		b.logger.Info("Failed to convert user ID ", err.Error())
	}
	userExists, err := b.storage.User().UserExists(m.From.Username)
	if err != nil {
		b.logger.Info("Failed to find user: ", err.Error())
	}
	user, ok, err := b.storage.User().FindByTelegramUsernameWithGreetings(m.From.Username)
	if err != nil {
		b.logger.Info("Failed to find user: ", err.Error())
	}
	if ok {
		msg = user.Greeting.Message
	} else if !ok && !userExists {
		err = b.storage.User().Create(&models.User{Username: m.From.Username, TelegramId: &userId})
		if err != nil {
			b.logger.Info("Failed to create new user: ", err.Error())
		}
		msg = stdMessage
	} else {
		msg = stdMessage
	}

	b.LogHandler(m, msg)
	handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
}
