package kafka

import (
	"checker/pkg/storage/esearch"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strings"
)

type ClientKafka struct {
	brokerAddress string
	port          string
	network       string
	eventTopic    string

	writer *kafka.Writer
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
	k.writer = &kafka.Writer{
		Addr:     kafka.TCP(k.brokerAddress),
		Topic:    k.eventTopic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (k *ClientKafka) ProduceEvents(ctx context.Context, docs []esearch.ElasticDocs) error {
	for _, doc := range docs {
		log.Println("[INFO] sending data of doc : ", doc.Artist)
		eventJson, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("[ERR] cannot unmarshal the event : %s\n", err.Error())
		}

		err = k.ProduceEvent(ctx, eventJson)
		if err != nil {
			return fmt.Errorf("[ERR] cannot send the event to kafka : %s", err.Error())
		}
	}
	return nil
}

func (k *ClientKafka) ProduceEvent(ctx context.Context, eventJson []byte) error {
	err := k.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte("artist-event"),
			Value: eventJson,
		})
	if err != nil {
		return fmt.Errorf("can't send info to kafka : %s", err.Error())
	}
	return nil
}
