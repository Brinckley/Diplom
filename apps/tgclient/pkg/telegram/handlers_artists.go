package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"tgclient/pkg/utils"
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
			if len(albums) <= 10 {
				albumsMsg := "Альбомы артиста " + name + " :\n" + utils.ListOfAlbumsFormatter(albums)
				msg := tgbotapi.NewMessage(message.Chat.ID, albumsMsg)
				_, err = b.bot.Send(msg)
				return ResOk, err
			} else {
				left := 10
				for i := 0; ; {
					albumsMsg := "Альбомы артиста " + name + " :\n" + utils.ListOfAlbumsFormatter(albums[i:left])
					msg := tgbotapi.NewMessage(message.Chat.ID, albumsMsg)
					_, err = b.bot.Send(msg)
					if left+10 < len(albums) {
						i = left
						left += 10
					} else {
						if left != len(albums) {
							i = left
							left = len(albums)
						} else {
							break
						}
					}
				}
			}
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
