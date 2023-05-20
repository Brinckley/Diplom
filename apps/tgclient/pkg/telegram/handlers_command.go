package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tgclient/pkg/kafka"
)

const (
	startCmd       = "start"
	helpCmd        = "help"
	artistsCmd     = "artists"
	albumsCmd      = "albums"
	favoritesCmd   = "favorites"
	subscribeCmd   = "subscribe"
	unsubscribeCmd = "unsubscribe"
	catalogCmd     = "catalog"
)

const (
	ResOk = iota
	ResFail
	ResWaitForRemove
	ResWaitForAdd
)

func (b *Bot) handleRandomMessage(message *tgbotapi.Message, n *int) error {
	b.logger.Printf("[%s] : '%s'", message.From.UserName, message.Text)
	*n--
	txt := msgRandomMessage
	if *n == 0 {
		txt = msgRandomMessageFinal
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, txt)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommand(message *tgbotapi.Message) (int, error) {

	switch message.Command() {
	case startCmd:
		return b.handleStartCmd(message)
	case catalogCmd:
		return b.handleCatalog(message)
	case artistsCmd:
		return b.handleAllArtists(message)
	case albumsCmd:
		return b.handleAllArtists(message)
	case favoritesCmd:
		return b.handleAllFavorites(message)
	case subscribeCmd:
		return b.handleSubscriptionIntro(message)
	case unsubscribeCmd:
		return b.handleUnsubscriptionIntro(message)
	case helpCmd:
		return b.handleHelpCmd(message)
	default:
		return b.handleArtistCheck(message)
	}
}

// getAlbumsByArtist postgres
func (b *Bot) handleCatalog(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Check the menu")
	msg.ReplyMarkup = b.menuKeyboard
	_, err := b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleNewEventReceived(chatId int64, event kafka.Event) (int, error) {
	txtMsg := event.CreateNotification()
	msg := tgbotapi.NewMessage(chatId, txtMsg)
	_, err := b.bot.Send(msg)
	return ResOk, err
}
