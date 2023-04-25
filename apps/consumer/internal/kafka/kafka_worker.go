package kafka

import (
	"consumer/internal/postgres"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
	"time"
)

type ClientKafka struct {
	cBrokerAddress string
	cNetwork       string
	cArtistTopic   string
	cAlbumTopic    string
	cTrackTopic    string

	cgArtist string
	cgAlbum  string
	cgTrack  string

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
	k.cgArtist = os.Getenv("ARTIST_CONSUMER_GROUP")
	k.cgAlbum = os.Getenv("ALBUM_CONSUMER_GROUP")
	k.cgTrack = os.Getenv("TRACK_CONSUMER_GROUP")
}

func (k *ClientKafka) ConsumeAndSend() {
	topicChan := make(chan string)
	var m sync.Mutex
	go k.consumeArtist(topicChan, &m)
	//	go k.consumeAlbum(topicChan)
	//	go k.consumeTrack(topicChan)

	timeout := time.After(70 * time.Second)
	for i := 0; i < 1; i++ {
		select {
		case tc := <-topicChan:
			log.Println("Topic Parsed :", tc)
		case <-timeout:
			log.Println("Error getting info from topic and writing to the db (timeout)!")
		}
	}
	//k.postgresClient.DBSelectArtists()
	//k.postgresClient.DBSelectAlbums()
	//k.postgresClient.DBSelectTracks()
}

func (k *ClientKafka) consumeArtist(tcs chan string, mutex *sync.Mutex) {
	//group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
	//	ID:                     "art-cg",
	//	Brokers:                []string{k.cBrokerAddress},
	//	Topics:                 []string{k.cArtistTopic},
	//	StartOffset:            kafka.LastOffset,
	//	ErrorLogger:            k.logger,
	//	Timeout:                time.Second * 30,
	//})
	c := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{k.cBrokerAddress},
		Topic:       k.cArtistTopic,
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
	})

	//err := c.Subscribe(k.cArtistTopic, nil)
	//if err != nil {
	//	k.logger.Fatal("[ERR] can't create consumer for topic Artist")
	//}

	for {
		m, err := c.FetchMessage(context.Background())
		if err != nil {
			break
		}

		var artist postgres.ArtistDB
		err = json.Unmarshal(m.Value, &artist)
		if err != nil {
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", k.cBrokerAddress, k.cArtistTopic)
			artist = postgres.ArtistDB{}
		}

		if err != nil {
			return
		}
		log.Println("Appended data of artist :", artist.Name)
		//k.postgresClient.DBInsertArtist(artist) // Taken away for while debugging // adding value to postgres
	}
	if err := c.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	tcs <- k.cArtistTopic
}

func (k *ClientKafka) consumeAlbum(tcs chan string) {
	c := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.cBrokerAddress},
		Topic:    k.cArtistTopic,
		GroupID:  "alb-gr-id",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := c.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var album postgres.AlbumDB
		err = json.Unmarshal(m.Value, &album)
		if err != nil {
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", k.cBrokerAddress, k.cAlbumTopic)
			album = postgres.AlbumDB{}
		}
		log.Println("Appended data of album :", album.Name)
		//k.postgresClient.DBInsertAlbum(album) // Taken away for while debugging
	}
	if err := c.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	tcs <- k.cAlbumTopic
}

func (k *ClientKafka) consumeTrack(tcs chan string) {
	c := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.cBrokerAddress},
		Topic:    k.cArtistTopic,
		GroupID:  "tra-gr-id",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := c.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var track postgres.TrackDB
		err = json.Unmarshal(m.Value, &track)
		if err != nil {
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", k.cBrokerAddress, k.cTrackTopic)
			track = postgres.TrackDB{}
		}
		log.Println("Appended data of track :", track.Name)
		//k.postgresClient.DBInsertTrack(track) // Taken away for while debugging
	}
	if err := c.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")

	tcs <- k.cTrackTopic
}
