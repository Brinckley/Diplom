package kafka

import (
	"consumer/internal/postgres"
	"context"
	"encoding/json"
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
	ctx := context.Background()
	go k.consumeArtist(topicChan, ctx)
	go k.consumeAlbum(topicChan, ctx)
	go k.consumeTrack(topicChan, ctx)

	timeout := time.After(70 * time.Second)
	for i := 0; i < 1; i++ {
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

func (k *ClientKafka) consumeArtist(tcs chan string, ctx context.Context) {
	dialer := &kafka.Dialer{
		Timeout:   time.Second * 10,
		DualStack: true,
	}

	c := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{k.cBrokerAddress},
		Topic:           k.cArtistTopic,
		StartOffset:     kafka.FirstOffset,
		MinBytes:        1,
		MaxBytes:        100e6,
		CommitInterval:  time.Second / 100,
		GroupID:         "art-cg",
		MaxWait:         time.Hour * 24,
		ReadLagInterval: 1 * time.Second,
		Dialer:          dialer,
		QueueCapacity:   100 * 2,
	})

	for {
		m, err := c.ReadMessage(ctx)
		if err != nil {
			break
		}
		err = c.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println("[ERR] can't commit msg : ", err.Error())
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
		k.postgresClient.DBInsertArtist(artist) // Taken away for while debugging // adding value to postgres
	}
	if err := c.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	tcs <- k.cArtistTopic
}

func (k *ClientKafka) consumeAlbum(tcs chan string, ctx context.Context) {
	dialer := &kafka.Dialer{
		Timeout:   time.Second * 10,
		DualStack: true,
	}

	c := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{k.cBrokerAddress},
		Topic:           k.cAlbumTopic,
		StartOffset:     kafka.FirstOffset,
		MinBytes:        1,
		MaxBytes:        100e6,
		CommitInterval:  time.Second / 100,
		GroupID:         "alb-group-id",
		MaxWait:         time.Hour * 24,
		ReadLagInterval: 1 * time.Second,
		Dialer:          dialer,
		QueueCapacity:   100 * 2,
	})

	for {
		m, err := c.ReadMessage(ctx)
		if err != nil {
			break
		}
		err = c.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println("[ERR] can't commit msg : ", err.Error())
			break
		}

		var album postgres.AlbumDB
		err = json.Unmarshal(m.Value, &album)
		if err != nil {
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", k.cBrokerAddress, k.cAlbumTopic)
			album = postgres.AlbumDB{}
		}
		log.Println("Appended data of album :", album.Name)
		k.postgresClient.DBInsertAlbum(album) // Taken away for while debugging
	}
	if err := c.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	tcs <- k.cAlbumTopic
}

func (k *ClientKafka) consumeTrack(tcs chan string, ctx context.Context) {
	dialer := &kafka.Dialer{
		Timeout:   time.Second * 10,
		DualStack: true,
	}

	c := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{k.cBrokerAddress},
		Topic:           k.cTrackTopic,
		StartOffset:     kafka.FirstOffset,
		MinBytes:        1,
		MaxBytes:        100e6,
		CommitInterval:  time.Second / 100,
		GroupID:         "tra-group-id",
		MaxWait:         time.Hour * 24,
		ReadLagInterval: 1 * time.Second,
		Dialer:          dialer,
		QueueCapacity:   100 * 2,
	})

	timerStart := time.Now()
	trackCounter := 0
	for {
		trackCounter++
		if trackCounter == 1000 {
			log.Println("[TIME INFO] time from start (or previous zero countdown mark) :", time.Since(timerStart))
			timerStart = time.Now()
		}
		m, err := c.ReadMessage(ctx)
		if err != nil {
			break
		}
		err = c.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println("[ERR] can't commit msg : ", err.Error())
			break
		}

		var track postgres.TrackDB
		err = json.Unmarshal(m.Value, &track)
		if err != nil {
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", k.cBrokerAddress, k.cTrackTopic)
			track = postgres.TrackDB{}
		}
		//log.Println("Appended data of track :", track.Name)
		k.postgresClient.DBInsertTrack(track) // Taken away for while debugging
	}
	if err := c.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")

	tcs <- k.cTrackTopic
}
