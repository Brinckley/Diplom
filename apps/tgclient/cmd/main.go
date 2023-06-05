package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
	"tgclient/pkg/kafka"
	"tgclient/pkg/prometheus"
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

	promClient := prometheus.NewClientPrometheus()

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
	//wg.Add(3)
	wg.Add(4)
	chanEvent := make(chan kafka.Event)

	tgBot := telegram.NewBot(bot, &logger, pStorage, chanEvent, promClient) // creating custom bot client
	clKafka := kafka.NewKafka()

	go promClient.StartHandling()
	go clKafka.ConsumeEvents(context.Background(), chanEvent, &wg)
	go tgBot.StartHandlingEventUpdates(&wg)
	go func(group *sync.WaitGroup) {
		err := tgBot.StartHandling(group)
		if err != nil {
			log.Fatalln("[ERR] Unable to start bot")

		}
	}(&wg)
	wg.Wait()
}
