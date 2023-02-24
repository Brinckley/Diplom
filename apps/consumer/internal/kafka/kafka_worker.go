package kafka

import (
	"consumer/internal/postgres"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

var cBrokerAddress = ""
var cNetwork = ""
var cArtistTopic = ""
var cAlbumTopic = ""
var cTrackTopic = ""

func InitConsumer() {
	cBrokerAddress = os.Getenv("BROKER_ADDRESS")
	cNetwork = os.Getenv("NETWORK")
	cArtistTopic = os.Getenv("ARTIST_TOPIC_NAME")
	cAlbumTopic = os.Getenv("ALBUM_TOPIC_NAME")
	cTrackTopic = os.Getenv("TRACK_TOPIC_NAME")
}

func ConsumeAndSend() {
	topicChan := make(chan string)
	go consumeArtist(cArtistTopic, topicChan)
	go consumeAlbum(cAlbumTopic, topicChan)
	go consumeTrack(cTrackTopic, topicChan)

	timeout := time.After(30 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case tc := <-topicChan:
			log.Println("Topic Parsed :", tc)
		case <-timeout:
			log.Println("Error getting info from topic and writing to the db (timeout)!")
		}
	}
	postgres.DBSelectArtists()
	postgres.DBSelectAlbums()
	postgres.DBSelectTracks()
}

func consumeArtist(topicArtist string, tcs chan string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cBrokerAddress},
		Topic:    topicArtist,
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
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", cBrokerAddress, topicArtist)
			artist = postgres.ArtistDB{}
		}
		log.Println("Appended data of artist :", artist.Name)
		postgres.DBInsertArtist(artist) // adding value to postgres
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	tcs <- topicArtist
}

func consumeAlbum(topicAlbum string, tcs chan string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cBrokerAddress},
		Topic:    topicAlbum,
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
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", cBrokerAddress, topicAlbum)
			album = postgres.AlbumDB{}
		}
		log.Println("Appended data of album :", album.Name)
		postgres.DBInsertAlbum(album)
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")
	tcs <- topicAlbum
}

func consumeTrack(topicTrack string, tcs chan string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cBrokerAddress},
		Topic:    topicTrack,
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
			log.Printf("Error unmarshalling data from Broker : %v, Topic : %v\n", cBrokerAddress, topicTrack)
			track = postgres.TrackDB{}
		}
		log.Println("Appended data of track :", track.Name)
		postgres.DBInsertTrack(track)
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")

	tcs <- topicTrack
}
