package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
	"tgclient/pkg/kafka"
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
	logger := logrus.Logger{Formatter: &logrus.JSONFormatter{}} // creating and initializing the logger
	logger.SetOutput(os.Stdout)

	pStorage := storage.NewPostgres() // initializing database
	logger.Println("Postgresql initialized")

	bot, err := tgbotapi.NewBotAPI(tgToken) // creating connection to tg bot
	if err != nil {
		log.Fatalln("[ERR] Can't receive token")
	}
	bot.Debug = true

	var wg sync.WaitGroup
	wg.Add(2)
	chanEvent := make(chan kafka.Event)

	tgBot := telegram.NewBot(bot, &logger, pStorage, chanEvent) // creating custom bot client
	if err := tgBot.Start(); err != nil {
		log.Fatalln("[ERR] Unable to start bot")
	}
}

// when update time?!

// getting info from kafka
// select users that are connected with this artist
// meanwhile a wait channel is working. It is waiting for getting info about artist
// if matching user found, sending userId and event info to the waiting channel
// the channel sends info to the user himself
