package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tgclient/pkg/storage"
)

func (b *Bot) handleSubscriptionIntro(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgSubscribeQuestion)
	_, err := b.bot.Send(msg)
	return ResWaitForAdd, err
}

func (b *Bot) handleSubscription(message *tgbotapi.Message) (int, error) {
	subscribe, err := b.storage.Subscribe(message)
	if err == storage.ErrNoArtists {
		msg := tgbotapi.NewMessage(message.Chat.ID, msgSubscribeFail)
		_, err = b.bot.Send(msg)
		return ResOk, err
	}
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

func (b *Bot) handleNoFavorites(message *tgbotapi.Message) (int, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, msgNoFavorites)
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
