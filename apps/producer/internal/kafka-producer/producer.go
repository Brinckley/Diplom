package kafka_producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"producer/internal/parsers/parser-discogs/discogs-functions"
	"producer/internal/parsers/parser-lastfm/lastfm-functions"
	"strings"
)

var cBrokerAddress = "localhost:9092"
var cNetwork = "tcp"
var cArtistTopic = "Artist"
var cAlbumTopic = "Album"
var cTrackTopic = "Track"
var KafkaPort = "9092"

func InitProducer() {
	cBrokerAddress = os.Getenv("BROKER_ADDRESS")
	KafkaPort = strings.Split(cBrokerAddress, ":")[1]
	cNetwork = os.Getenv("NETWORK")
	cArtistTopic = os.Getenv("ARTIST_TOPIC_NAME")
	cAlbumTopic = os.Getenv("ALBUM_TOPIC_NAME")
	cTrackTopic = os.Getenv("TRACK_TOPIC_NAME")
	lastfm_functions.CommonRoot = os.Getenv("COMMON_ROOT")
	discogs_functions.CommonRoot = os.Getenv("COMMON_ROOT")
}

func Produce(topic string, message []byte, cs chan string) {
	switch topic {
	case "Artist":
		topic = cArtistTopic
	case "Album":
		topic = cAlbumTopic
	case "Track":
		topic = cTrackTopic
	default:
		return
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(cBrokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	messages := kafka.Message{
		Key:   []byte("Key-tmp"),
		Value: message,
	}
	err := w.WriteMessages(context.Background(), messages)
	if err != nil {
		panic("DSDs")
	}
	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	if cs != nil {
		cs <- topic
	}
}
