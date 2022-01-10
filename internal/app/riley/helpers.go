package riley

import (
	"github.com/drTragger/rileyBot/storage"
	"github.com/sirupsen/logrus"
)

func (b *Bot) configureLoggerField() error {
	logLevel, err := logrus.ParseLevel(b.config.LoggerLevel)
	if err != nil {
		return err
	}
	b.logger.SetLevel(logLevel)
	return nil
}

func (b *Bot) configureRouterField() {
	b.riley.HandleMessage("/start", b.StartHandler)
	b.riley.HandleMessage("/play", b.PlayHandler)
	b.riley.HandleMessage("/weather", b.cityRequestHandler)
	b.riley.HandleMessage("", b.weatherHandler)
	b.riley.HandleCallback(b.CallbackHandler)
}

func (b *Bot) configureStorageField() error {
	newStorage := storage.New(b.config.Storage)
	if err := newStorage.Open(); err != nil {
		return err
	}
	b.storage = newStorage
	return nil
}
