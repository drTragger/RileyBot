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
	translations = map[string]string{"scissors": "ножиці✂", "rock": "камінь🗿", "paper": "папір\U0001F9FB"}
)

func (b *Bot) PlayHandler(m *tbot.Message) {
	handleChatActionError(b.client.SendChatAction(m.Chat.ID, tbot.ActionTyping))
	time.Sleep(500 * time.Millisecond)
	buttons := makeButtons()

	b.LogHandler(m, "Showed buttons")
	handleMessageError(b.client.SendMessage(m.Chat.ID, "Твій ход:", tbot.OptInlineKeyboardMarkup(buttons)))
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
		Text:         "Камінь",
		CallbackData: "rock",
	}
	btnPaper := tbot.InlineKeyboardButton{
		Text:         "Папір",
		CallbackData: "paper",
	}
	btnScissors := tbot.InlineKeyboardButton{
		Text:         "Ножиці",
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
		result = "Нічия"
	case options[botMove]:
		result = "Ти програв"
	default:
		result = "Ти переміг"
	}
	msg = fmt.Sprintf("%s!\nТи обрав %s\nЯ обрав %s", result, translations[humanMove], translations[botMove])
	return
}
