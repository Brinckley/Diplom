package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func (b *Bot) handleArtistCheck(message *tgbotapi.Message) (int, error) {
	if len(message.Text) <= 1 {
		return b.handleUnknownCmd(message)
	}

	names, err := b.storage.GetAllArtists()
	if err != nil {
		return ResFail, err
	}

	for _, name := range names {
		nameOld := strings.Replace(name, " ", "", -1)
		if nameOld == message.Text[1:] {
			albums, err := b.storage.GetAlbumsByArtist(name)
			if err != nil {
				return 0, err
			}

			//albumsMsg := "Albums of artist " + name + " :\n" + albums
			albumsMsg := "Альбомы артиста " + name + " :\n" + albums
			msg := tgbotapi.NewMessage(message.Chat.ID, albumsMsg)
			_, err = b.bot.Send(msg)
			return ResOk, err
		}
	}

	return b.handleUnknownCmd(message)
}

func (b *Bot) handleAllArtists(message *tgbotapi.Message) (int, error) {
	//msg := tgbotapi.NewMessage(message.Chat.ID, "Check the menu")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Смотри меню")
	msg.ReplyMarkup = b.artistsKeyboard
	_, err := b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleNoArtists(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgNoArtists)
	_, err := b.bot.Send(msg)
	return err
}
