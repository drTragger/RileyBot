package riley

import (
	"github.com/drTragger/rileyBot/storage"
	"github.com/sirupsen/logrus"
	"github.com/yanzay/tbot/v2"
	"time"
)

type Bot struct {
	client  *tbot.Client
	riley   *tbot.Server
	config  *Config
	logger  *logrus.Logger
	storage *storage.Storage
}

func New(config *Config, client *tbot.Client, bot *tbot.Server) *Bot {
	return &Bot{
		config: config,
		logger: logrus.New(),
		client: client,
		riley:  bot,
	}
}

func (b *Bot) Start() error {
	if err := b.configureLoggerField(); err != nil {
		return err
	}
	b.logger.Info("Started Riley riley")
	b.configureRouterField()
	if err := b.configureStorageField(); err != nil {
		return err
	}
	return b.riley.Start()
}

func (b *Bot) LogHandler(m *tbot.Message) {
	b.logger.Printf("%s\nUsername: %s\nChat ID: %s\nMessage: %s\n============================", time.Now().UTC(), m.From.Username, m.Chat.ID, m.Text)
}

func (b *Bot) LogCallbackHandler(cq *tbot.CallbackQuery) {
	b.logger.Printf("%s\nUsername: %s\nChat ID: %s\nMessage: %s\n============================", time.Now().UTC(), cq.From.Username, cq.Message.Chat.ID, cq.Message.Text)
}
