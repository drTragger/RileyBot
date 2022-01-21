package riley

import (
	"fmt"
	"github.com/yanzay/tbot/v2"
	"time"
)

func (b *Bot) jokesHandler(m *tbot.Message) {
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)

	var msg string
	joke, ok, err := b.storage.Joke().GetRandom()
	if err != nil {
		b.logger.Info("Error during fetching joke data: ", err.Error())
		msg = "Извините, временно туплю.\nПожалуйста, попробуйте позже.\nА пока можете поиграть в Камень-Ножницы-Бумага /play"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}
	if !ok {
		b.logger.Info("No jokes found")
		msg = "Извините, мой дед не рассказывал мне анекдоты.\nПожалуйста, попробуйте позже.\nЯ схожу к деду за порцией анекдотов.\nА пока можете поиграть в Камень-Ножницы-Бумага /play"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	fmt.Println(joke.Joke)
	msg = fmt.Sprintf("Внимание, анекдот!\n\n%s", joke.Joke)

	b.LogHandler(m, msg)
	handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
}
