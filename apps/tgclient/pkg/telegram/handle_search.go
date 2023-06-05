package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (b *Bot) handleSearch(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgSearchArtist)
	_, err := b.bot.Send(msg)
	return ResWaitForAdd, err
}
