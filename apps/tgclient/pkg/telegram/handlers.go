package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tgclient/pkg/storage"
)

const (
	startCmd       = "start"
	helpCmd        = "help"
	artistsCmd     = "artists"
	favoritesCmd   = "favorites"
	subscribeCmd   = "subscribe"
	unsubscribeCmd = "unsubscribe"
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
	// THINK : return flag of waiting for new msg?
	b.logger.Printf("[%s] : '%s'", message.From.UserName, message.Text)

	switch message.Command() {
	case startCmd:
		return b.handleStartCmd(message)
	case artistsCmd:
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
		return b.handleUnknownCmd(message)
	}
}

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

func (b *Bot) handleAllArtists(message *tgbotapi.Message) (int, error) {
	artists, err := b.storage.GetAllArtists()
	if err == storage.ErrNoArtists {
		return ResFail, b.handleNoArtists(message)
	}
	if err != nil {
		return ResFail, err
	}
	var all string
	for _, artist := range artists {
		all += artist + "\n"
	}
	if len(all) == 0 {
		all = msgNoArtists
	} else {
		all = msgIntroArtists + all
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, all)
	_, err = b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleNoArtists(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgNoArtists)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleAllFavorites(message *tgbotapi.Message) (int, error) {
	favorites, err := b.storage.GetFavorites(message)
	if err == storage.ErrNoFavorites {
		return b.handleNoFavorites(message)
	} else if err != nil {
		return ResFail, err
	}
	var all string
	for _, f := range favorites {
		all += f + "\n"
	}
	if len(all) == 0 {
		all = msgNoFavorites
	} else {
		all = msgIntroFavorites + all
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, all)
	_, err = b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleSubscriptionIntro(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgSubscribeQuestion)
	_, err := b.bot.Send(msg)
	return ResWaitForAdd, err
}

func (b *Bot) handleSubscription(message *tgbotapi.Message) (int, error) {
	subscribe, err := b.storage.Subscribe(message)
	if err != nil {
		return ResFail, err
	}
	if subscribe {
		msg := tgbotapi.NewMessage(message.Chat.ID, msgSubscribeSuccess)
		_, err = b.bot.Send(msg)
		return ResOk, err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgSubscribeAlready)
	_, err = b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleUnsubscriptionIntro(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgUnsubscribeQuestion)
	_, err := b.bot.Send(msg)
	return ResWaitForRemove, err
}

func (b *Bot) handleUnsubscription(message *tgbotapi.Message) (int, error) {
	subscribe, err := b.storage.Unsubscribe(message)
	if err != nil {
		return ResFail, err
	}
	if subscribe {
		msg := tgbotapi.NewMessage(message.Chat.ID, msgUnsubscribeSuccess)
		_, err = b.bot.Send(msg)
		return ResOk, err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgUnsubscribeFail)
	_, err = b.bot.Send(msg)
	return ResOk, nil
}

func (b *Bot) handleNoFavorites(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgNoFavorites)
	_, err := b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleUnknownCmd(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgUnknownCommand)
	_, err := b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleHelpCmd(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgHelpCommand)
	_, err := b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleNotArtistNameSub(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgSubscribeFail)
	_, err := b.bot.Send(msg)
	return ResOk, err
}

func (b *Bot) handleNotArtistNameUnsub(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgUnsubscribeFail)
	_, err := b.bot.Send(msg)
	return ResOk, err
}
