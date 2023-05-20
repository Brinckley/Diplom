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

			albumsMsg := "Albums of artist " + name + " :\n" + albums
			msg := tgbotapi.NewMessage(message.Chat.ID, albumsMsg)
			_, err = b.bot.Send(msg)
			return ResOk, err
		}
	}

	return b.handleUnknownCmd(message)
}

func (b *Bot) handleAllArtists(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Check the menu")
	msg.ReplyMarkup = b.artistsKeyboard
	_, err := b.bot.Send(msg)
	return ResOk, err
	//artists, err := b.storage.GetAllArtists()
	//if err == storage.ErrNoArtists {
	//	return ResFail, b.handleNoArtists(message)
	//}
	//if err != nil {
	//	return ResFail, err
	//}
	//var all string
	//for _, artist := range artists {
	//	all += artist + "\n"
	//}
	//if len(all) == 0 {
	//	all = msgNoArtists
	//} else {
	//	all = msgIntroArtists + all
	//}
	//msg := tgbotapi.NewMessage(message.Chat.ID, all)
	//_, err = b.bot.Send(msg)
	//return ResOk, err
}

func (b *Bot) handleNoArtists(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgNoArtists)
	_, err := b.bot.Send(msg)
	return err
}
