package riley

import (
	"fmt"
	"github.com/yanzay/tbot/v2"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	picks   = []string{"rock", "paper", "scissors"}
	options = map[string]string{
		"paper":    "rock",
		"scissors": "paper",
		"rock":     "scissors",
	}
	translations = map[string]string{"scissors": "ножницы✂", "rock": "камень🗿", "paper": "бумагу\U0001F9FB"}
)

func (b *Bot) PlayHandler(m *tbot.Message) {
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)
	buttons := makeButtons()

	b.LogHandler(m, "Showed buttons")
	handleMessageError(b.client.SendMessage(m.Chat.ID, "Твой ход:", tbot.OptInlineKeyboardMarkup(buttons)))
}

func (b *Bot) CallbackHandler(cq *tbot.CallbackQuery) {
	handleChatActionError(b.client.SendChatAction(cq.Message.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)
	humanMove := cq.Data
	msg := playGame(humanMove)

	b.LogCallbackHandler(cq, msg)
	handleChatActionError(b.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID))
	handleMessageError(b.client.SendMessage(cq.Message.Chat.ID, msg))
}

func makeButtons() *tbot.InlineKeyboardMarkup {
	btnRock := tbot.InlineKeyboardButton{
		Text:         "Камень",
		CallbackData: "rock",
	}
	btnPaper := tbot.InlineKeyboardButton{
		Text:         "Бумага",
		CallbackData: "paper",
	}
	btnScissors := tbot.InlineKeyboardButton{
		Text:         "Ножницы",
		CallbackData: "scissors",
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{btnRock, btnScissors, btnPaper},
		},
	}
}

func playGame(humanMove string) (msg string) {
	var result string
	botMove := picks[rand.Intn(len(picks))]
	switch humanMove {
	case botMove:
		result = "Ничья"
	case options[botMove]:
		result = "Ты проиграл"
	default:
		result = "Ты выиграл"
	}
	msg = fmt.Sprintf("%s!\nТы выбрал %s\nЯ выбрал %s", result, translations[humanMove], translations[botMove])
	return
}
