package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"tgclient/pkg/kafka"
	"tgclient/pkg/storage"
	"tgclient/pkg/utils"
	"time"
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

func (b *Bot) setKeyboards() {
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

	b.artistsKeyboard = tgbotapi.NewInlineKeyboardMarkup()
}

func (b *Bot) initUpdatesChannel() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	b.setKeyboards()
	var err error
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

func (b *Bot) StartHandlingEventUpdates(group *sync.WaitGroup) {
	defer group.Done()
	for {
		select {
		case event, ok := <-b.receiverChan:
			if ok {
				log.Println("[INFO] event received")
				subscribers, err := b.storage.GetAllSubscribers(event)
				if err != nil {
					log.Println("[ERR] can't receive event from kafka!")
					return
				}

				log.Println("[INFO] starting iteration over list of subscribers...")
				log.Println("[INFO] all list of subs : ", subscribers)

				for _, sub := range subscribers {
					log.Printf("[INFO STEP] int64 val : %v\n", sub)
					_, err := b.handleNewEventReceived(sub, event)
					if err != nil {
						log.Printf("[ERR] can't handle event from kafka for chat %v!\n", sub)
						return
					}
				}
			}
		default:
			b.logger.Println("[UPD] no event updates found")
		}
		time.Sleep(10 * time.Second)
	}
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	prev := "none"
	rand.Seed(time.Now().UnixNano())
	counter := rand.Intn(10) + 4
	for update := range updates {
		select {
		case event, ok := <-b.receiverChan:
			if ok {
				log.Println("[INFO] event received")
				subscribers, err := b.storage.GetAllSubscribers(event)
				if err != nil {
					log.Println("[ERR] can't receive event from kafka!")
					return
				}
				log.Println("[INFO] starting iteration over list of subscribers...")
				log.Println("[INFO] all list of subs : ", subscribers)
				for _, sub := range subscribers {
					log.Printf("[INFO STEP] int64 val : %v\n", sub)
					_, err := b.handleNewEventReceived(sub, event)
					if err != nil {
						log.Printf("[ERR] can't handle event from kafka for chat %v!\n", sub)
						return
					}
				}
			}
		default:
			b.logger.Println("[UPD] no event updates found")
		}
		fmt.Println("Prev cmd : ", prev)

		if prev == "/subscribe" {
			b.waitForAdd(update)
			prev = update.Message.Text
			continue
		}

		if prev == "/unsubscribe" {
			b.waitForRemove(update)
			prev = update.Message.Text
			continue
		}

		if update.Message != nil {
			if update.Message.IsCommand() {
				prev = update.Message.Text
				cmdRes, err := b.handleCommand(update.Message)
				if err != nil {
					b.logger.Printf("[ERR] Can't handle command '%s', err : %s", update.Message.Text, err.Error())
					continue
				}
				b.logger.Printf("After handling command '%s' got result : %s", update.Message.Text, strconv.Itoa(cmdRes))
				continue
			} else {
				err := b.handleRandomMessage(update.Message, &counter)
				if counter == 0 {
					counter = rand.Intn(10) + 4
				}
				if err != nil {
					b.logger.Printf("[ERR] Can't handle message '%s', err : %s", update.Message.Text, err.Error())
				}
			}
		} else if update.CallbackQuery != nil {
			m, err := b.callbackHandler(&update)
			if err != nil {
				log.Println(err.Error())
			}
			prev = "/" + m
			continue
		}
	}
	time.Sleep(5 * time.Second)
}

func (b *Bot) callbackHandler(update *tgbotapi.Update) (string, error) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := b.bot.AnswerCallbackQuery(callback); err != nil {
		return "", fmt.Errorf("[ERR] unable to create callback query : %s", err.Error())
	}

	m := utils.CallbackToMsg(update.CallbackQuery)
	_, err := b.handleCommand(&m)
	if err != nil {
		return "", fmt.Errorf("[ERR] with handling callback : %s", err.Error())
	}

	return m.Command(), nil
}

func (b *Bot) waitForAdd(update tgbotapi.Update) bool {
	if update.Message.IsCommand() {
		_, err := b.handleNotArtistNameSub(update.Message)
		if err != nil {
			b.logger.Printf("[ERR] Can't handle command '%s', err : %s", update.Message.Text, err.Error())
		}
		return false
	}

	_, err := b.handleSubscription(update.Message)
	if err != nil {
		b.logger.Printf("[ERR] Can't handle command '%s', err : %s", update.Message.Text, err.Error())
	}
	return true
}

func (b *Bot) waitForRemove(update tgbotapi.Update) bool {
	if update.Message.IsCommand() {
		_, err := b.handleNotArtistNameUnsub(update.Message)
		if err != nil {
			b.logger.Printf("[ERR] Can't handle command '%s', err : %s", update.Message.Text, err.Error())
		}
		return false
	}

	_, err := b.handleUnsubscription(update.Message)
	if err != nil {
		b.logger.Printf("[ERR] Can't handle command '%s', err : %s", update.Message.Text, err.Error())
	}
	return true
}
