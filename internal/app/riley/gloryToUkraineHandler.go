package riley

import (
	"github.com/yanzay/tbot/v2"
	"time"
)

func (b *Bot) GloryToUkraineHandler(m *tbot.Message) {
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)
	msg := "Героям Слава!\nСлава Нації!\nПиздець російській федерації!"
	b.LogHandler(m, msg)
	handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
}
