package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tgclient/pkg/storage"
)

func (b *Bot) handleStartCmd(message *tgbotapi.Message) (int, error) {
	err := b.storage.Registration(message)
	if err == storage.ErrUserExists {
		return b.handleStartAgainCmd(message)
	} else if err != nil {
		return ResFail, err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgStartCommand)
	_, err = b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleStartAgainCmd(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgStartAgainCommand)
	_, err := b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleHelpCmd(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgHelpCommand)
	_, err := b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleUnknownCmd(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgUnknownCommand)
	_, err := b.bot.Send(msg)
	return ResOk, err
}
