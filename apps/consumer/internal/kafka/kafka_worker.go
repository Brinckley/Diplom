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

	logger         *logrus.Logger
	postgresClient *postgres.ClientPostgres
}

func NewKafka(pc *postgres.ClientPostgres, logger *logrus.Logger) *ClientKafka {
	var k ClientKafka
	k.logger = logger
	k.postgresClient = pc
	k.init()
	return &k
}

func (k *ClientKafka) init() {
	k.cBrokerAddress = os.Getenv("BROKER_ADDRESS")
	k.cNetwork = os.Getenv("NETWORK")
	k.cArtistTopic = os.Getenv("ARTIST_TOPIC_NAME")
	k.cAlbumTopic = os.Getenv("ALBUM_TOPIC_NAME")
	k.cTrackTopic = os.Getenv("TRACK_TOPIC_NAME")
}

func (k *ClientKafka) ConsumeAndSend() {
	topicChan := make(chan string)
	go k.consumeArtist(topicChan)
	go k.consumeAlbum(topicChan)
	go k.consumeTrack(topicChan)

	timeout := time.After(30 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case tc := <-topicChan:
			log.Println("Topic Parsed :", tc)
		case <-timeout:
			log.Println("Error getting info from topic and writing to the db (timeout)!")
		}
	}
	k.postgresClient.DBSelectArtists()
	k.postgresClient.DBSelectAlbums()
	k.postgresClient.DBSelectTracks()
}

func (k *ClientKafka) consumeArtist(tcs chan string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.cBrokerAddress},
		Topic:    k.cArtistTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
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
		k.postgresClient.DBInsertArtist(artist) // adding value to postgres
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	tcs <- k.cArtistTopic
}

func (k *ClientKafka) consumeAlbum(tcs chan string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.cBrokerAddress},
		Topic:    k.cAlbumTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
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
		k.postgresClient.DBInsertAlbum(album)
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	tcs <- k.cAlbumTopic
}

func (k *ClientKafka) consumeTrack(tcs chan string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.cBrokerAddress},
		Topic:    k.cTrackTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
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
		k.postgresClient.DBInsertTrack(track)
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")

	tcs <- k.cTrackTopic
}
