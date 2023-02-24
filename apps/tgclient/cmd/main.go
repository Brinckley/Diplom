package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"tgclient/pkg/storage"
	"tgclient/pkg/telegram"
)

var (
	tgToken = ""
)

func InitTelegram() {
	tgToken = os.Getenv("TG_BOT_TOKEN")
}

func main() {
	InitTelegram()
	logger := logrus.Logger{Formatter: &logrus.JSONFormatter{}}
	logger.SetOutput(os.Stdout)

	storage := storage.NewPostgres()
	logger.Println("Postgresql initialized")

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Fatalln("[ERR] Can't receive token")
	}
	bot.Debug = true

	tgBot := telegram.NewBot(bot, &logger, storage)
	if err := tgBot.Start(); err != nil {
		log.Fatalln("[ERR] Unable to start bot")
	}
}
