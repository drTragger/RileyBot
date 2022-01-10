package riley

import (
	"github.com/yanzay/tbot/v2"
	"log"
)

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
