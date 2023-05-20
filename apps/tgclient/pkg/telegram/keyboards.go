package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func (b *Bot) setKeyboards() error {
	b.setKeyboardMainMenu()
	names, err := b.storage.GetAllArtists()
	if err != nil {
		return err
	}
	b.setKeyboardArtists(names)
	return nil
}

func (b *Bot) setKeyboardMainMenu() {
	b.menuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Artists", "/artists"),
			tgbotapi.NewInlineKeyboardButtonData("Albums", "/albums"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Subscribe", "/subscribe"),
			tgbotapi.NewInlineKeyboardButtonData("Unsubscribe", "/unsubscribe"),
			tgbotapi.NewInlineKeyboardButtonData("Favorites", "/favorites"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Help", "/help"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Discogs", "https://www.discogs.com/"),
			tgbotapi.NewInlineKeyboardButtonURL("Last.fm", "https://www.last.fm/ru/"),
			tgbotapi.NewInlineKeyboardButtonURL("Kassir.ru", "https://www.kassir.ru/"),
		),
	)
}

func (b *Bot) setKeyboardArtists(names []string) {
	b.artistsKeyboard = tgbotapi.NewInlineKeyboardMarkup()
	for i := 0; i < len(names); i += 3 {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(names[i], "/"+strings.Replace(names[i], " ", "", -1)))
		if i+1 < len(names) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(names[i+1], "/"+strings.Replace(names[i+1], " ", "", -1)))
		}
		if i+2 < len(names) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(names[i+2], "/"+strings.Replace(names[i+2], " ", "", -1)))
		}
		b.artistsKeyboard.InlineKeyboard = append(b.artistsKeyboard.InlineKeyboard, row)
	}
}
