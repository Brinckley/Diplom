package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type ClientKafka struct {
	brokerAddress string
	port          string
	network       string
	eventTopic    string
	l             *log.Logger

	reader *kafka.Reader
}

func NewKafka() *ClientKafka {
	var k ClientKafka
	k.init()
	return &k
}

func (k *ClientKafka) init() {
	k.brokerAddress = os.Getenv("BROKER_ADDRESS")
	k.port = strings.Split(k.brokerAddress, ":")[1]
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

		// sending msg data to channel
		ueChan <- event
	}
	wg.Done()
}
