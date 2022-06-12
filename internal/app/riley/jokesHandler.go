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
		msg = "Вибачте, тимчасово туплю.\nБудь ласка, спробуйте пізніше.\nА поки можете пограти у Камінь-Ножиці-Папір /play"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}
	if !ok {
		b.logger.Info("No jokes found")
		msg = "Вибачте, мій дід не розповідав мені анекдоти.\nБудь ласка, спробуйте пізніше.\nЯ піду до діда за порцією анектодів.\nА поки можете пограти у Камінь-Ножиці-Папір /play"
		handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
		return
	}

	fmt.Println(joke.Joke)
	msg = fmt.Sprintf("Увага, анекдот!\n\n%s", joke.Joke)

	b.LogHandler(m, msg)
	handleMessageError(b.client.SendMessage(m.Chat.ID, msg))
}
