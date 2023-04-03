package kafka

import (
	"consumer/internal/postgres"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

type ClientKafka struct {
	cBrokerAddress string
	cNetwork       string
	cArtistTopic   string
	cAlbumTopic    string
	cTrackTopic    string

	dbClient     *postgres.ClientPostgres
	readerArtist *kafka.Reader
	readerAlbum  *kafka.Reader
	readerTrack  *kafka.Reader
	topicChan    chan string

	logger *logrus.Logger
}

func NewKafka(db *postgres.ClientPostgres, logger *logrus.Logger) *ClientKafka {
	var k ClientKafka
	k.logger = logger
	k.dbClient = db
	k.init()
	return &k
}

func (k *ClientKafka) init() {
	k.cBrokerAddress = os.Getenv("BROKER_ADDRESS")
	k.cNetwork = os.Getenv("NETWORK")
	k.cArtistTopic = os.Getenv("ARTIST_TOPIC_NAME")
	k.cAlbumTopic = os.Getenv("ALBUM_TOPIC_NAME")
	k.cTrackTopic = os.Getenv("TRACK_TOPIC_NAME")
	k.topicChan = make(chan string)
	k.readerArtist = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{k.cBrokerAddress},
		Topic:       k.cArtistTopic,
		GroupID:     "artist-consumer-group",
		MinBytes:    5,
		MaxBytes:    1e6,
		StartOffset: kafka.LastOffset,
		MaxWait:     3 * time.Second,
		Logger:      k.logger,
	})
	k.readerAlbum = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{k.cBrokerAddress},
		Topic:       k.cAlbumTopic,
		GroupID:     "album-consumer-group",
		MinBytes:    5,
		MaxBytes:    1e6,
		StartOffset: kafka.LastOffset,
		MaxWait:     3 * time.Second,
		Logger:      k.logger,
	})
	k.readerArtist = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{k.cBrokerAddress},
		Topic:       k.cAlbumTopic,
		GroupID:     "track-consumer-group",
		MinBytes:    5,
		MaxBytes:    1e6,
		StartOffset: kafka.LastOffset,
		MaxWait:     3 * time.Second,
		Logger:      k.logger,
	})
}

func (k *ClientKafka) ConsumeAndSend() {
	go k.consumeArtist()
	go k.consumeAlbum()
	go k.consumeTrack()

	timeout := time.After(30 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case tc := <-k.topicChan:
			log.Println("Topic Parsed :", tc)
		case <-timeout:
			log.Println("Error getting info from topic and writing to the db (timeout)!")
		}
	}

	k.dbClient.DBSelectArtists()
	k.dbClient.DBSelectAlbums()
	k.dbClient.DBSelectTracks()
}

func (k *ClientKafka) consumeArtist() {
	for {
		m, err := k.readerArtist.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		var artist postgres.ArtistDB
		err = json.Unmarshal(m.Value, &artist)
		if err != nil {
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", k.cBrokerAddress, k.cArtistTopic)
			artist = postgres.ArtistDB{}
		}
		log.Println("Appended data of artist :", artist.Name)
		k.dbClient.DBInsertArtist(artist) // adding value to postgres
	}
	if err := k.readerArtist.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	k.topicChan <- k.cArtistTopic
}

func (k *ClientKafka) consumeAlbum() {
	for {
		m, err := k.readerAlbum.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		var album postgres.AlbumDB
		err = json.Unmarshal(m.Value, &album)
		if err != nil {
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", k.cBrokerAddress, k.cAlbumTopic)
			album = postgres.AlbumDB{}
		}
		log.Println("Appended data of album :", album.Name)
		k.dbClient.DBInsertAlbum(album)
	}
	if err := k.readerAlbum.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	k.topicChan <- k.cAlbumTopic
}

func (k *ClientKafka) consumeTrack() {
	for {
		m, err := k.readerTrack.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		var track postgres.TrackDB
		err = json.Unmarshal(m.Value, &track)
		if err != nil {
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", k.cBrokerAddress, k.cTrackTopic)
			track = postgres.TrackDB{}
		}
		log.Println("Appended data of track :", track.Name)
		k.dbClient.DBInsertTrack(track)
	}
	if err := k.readerTrack.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")

	k.topicChan <- k.cTrackTopic
}
