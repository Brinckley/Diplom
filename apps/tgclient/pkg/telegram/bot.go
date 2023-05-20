package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"sync"
	"tgclient/pkg/kafka"
	"tgclient/pkg/storage"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	logger       *logrus.Logger
	storage      *storage.TgPostgres
	receiverChan chan kafka.Event

	updates tgbotapi.UpdatesChannel

	menuKeyboard    tgbotapi.InlineKeyboardMarkup
	artistsKeyboard tgbotapi.InlineKeyboardMarkup
}

func NewBot(bot *tgbotapi.BotAPI, logger *logrus.Logger, storage *storage.TgPostgres, rc chan kafka.Event) *Bot {
	b := &Bot{
		bot:          bot,
		logger:       logger,
		storage:      storage,
		receiverChan: rc,
	}
	err := b.initUpdatesChannel()
	if err != nil {
		return nil
	}
	return b
}

func (b *Bot) initUpdatesChannel() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	err := b.setKeyboards()
	if err != nil {
		return err
	}
	b.updates, err = b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) StartHandling(wg *sync.WaitGroup) error {
	defer wg.Done()
	b.logger.Println("Authorised with username : ", b.bot.Self.UserName)
	for {
		b.handleUpdates(b.updates)
	}
	return nil
}
