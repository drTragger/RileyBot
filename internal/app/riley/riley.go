package riley

import (
	"github.com/drTragger/RileyBot/storage"
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

func (b *Bot) LogHandler(m *tbot.Message, answer string) {
	b.logger.Printf("%s\nUsername: %s\nChat ID: %s\nMessage: %s\nAnswer: %s\n============================", time.Now().UTC(), m.From.Username, m.Chat.ID, m.Text, answer)
}

func (b *Bot) LogCallbackHandler(cq *tbot.CallbackQuery, answer string) {
	b.logger.Printf("%s\nUsername: %s\nChat ID: %s\nMessage: %s\nnAnswer: %s\n============================", time.Now().UTC(), cq.From.Username, cq.Message.Chat.ID, cq.Message.Text, answer)
}
