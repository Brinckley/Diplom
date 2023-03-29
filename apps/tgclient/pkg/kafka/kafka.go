package kafka

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"sync"
	"time"
)

type ClientKafka struct {
	brokerAddress string
	port          string
	network       string
	eventTopic    string
	l             *log.Logger

	reader     *kafka.Reader
	tgUser     *tgbotapi.User
	senderChan *tgbotapi.UpdatesChannel
}

func NewKafka(channel *tgbotapi.UpdatesChannel) *ClientKafka {
	var k ClientKafka
	k.init()
	k.senderChan = channel
	return &k
}

func (k *ClientKafka) init() {
	log.Println("[INFO] kafka initialization started")
	k.brokerAddress = os.Getenv("BROKER_ADDRESS")
	k.network = os.Getenv("NETWORK")
	k.eventTopic = os.Getenv("EVENT_TOPIC_NAME")
	k.l = log.New(os.Stdout, "kafka event reader: ", 0)
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{k.brokerAddress},
		Topic:       k.eventTopic,
		MinBytes:    5,
		MaxBytes:    1e6,
		StartOffset: kafka.FirstOffset,
		MaxWait:     3 * time.Second,
		Logger:      k.l,
	})
	k.tgUser = &tgbotapi.User{
		ID:        -1,
		FirstName: "Kafka",
	}
	log.Println("[INFO] kafka initialization finished")
}

func (k *ClientKafka) ConsumeEvents(ctx context.Context, ueChan chan Event, wg *sync.WaitGroup) {
	for {
		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("[ERR] can't receive data from kafka event topic :", err.Error())
		}

		var event Event
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("[ERR] can't unmarshall data from Event : %s", err.Error())
			event = Event{}
			continue
		}

		time.Sleep(5 * time.Second)
		// sending msg data to channel
		log.Println("[INFO] received msg about artist : ", event.Artist)
		updEvent, err := k.createUpdateFromEvent(event)
		ueChan <- updEvent
	}
	wg.Done()
}

func (k *ClientKafka) createUpdateFromEvent(event Event) (tgbotapi.Update, error) {
	return tgbotapi.Update{
		UpdateID: -1,
		Message: &tgbotapi.Message{
			MessageID: -1,
			From:      k.tgUser,
			Date:      0,
			Text:      event.CreateNotification(),
		},
	}, nil
}
