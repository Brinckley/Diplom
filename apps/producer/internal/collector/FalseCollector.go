package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	kafka_producer "producer/internal/kafka-producer"
)

func FalseCollector(num int) {
	for i := 0; i < num; i++ {
		log.Println("[INFO] Sent data with ID: ", i)
		tr := TrackDB{
			Id:         i,
			Name:       fmt.Sprintf("Track%v", i),
			UrlLastfm:  "Tr",
			Duration:   "Tr",
			Position:   "Tr",
			ArtistHash: 0,
			AlbumHash:  0,
		}
		al := AlbumDB{
			Id:         i,
			Name:       fmt.Sprintf("Album%v", i),
			Release:    "",
			UrlLastfm:  "",
			UrlDiscogs: "",
			Picture:    "",
			TrackCount: i,
			ArtistHash: 0,
			AlbumHash:  0,
		}
		ar := ArtistDB{
			Id:         i,
			Name:       fmt.Sprintf("Artist%v", i),
			Bio:        "",
			OnTour:     false,
			Picture:    "",
			UrlLastfm:  "",
			UrlDiscogs: "",
			Genre:      "",
			IdLastfm:   "",
			IdDiscogs:  "",
			ArtistHash: 0,
		}

		trackJson, err := json.Marshal(tr)
		if err != nil {
			log.Println("Error marshalling track :", tr.Name)
		}
		trackJson = bytes.Trim(trackJson, "\x00")

		albumJson, err := json.Marshal(al)
		if err != nil {
			log.Println("Error marshalling track :", al.Name)
		}
		albumJson = bytes.Trim(albumJson, "\x00")

		artistJson, err := json.Marshal(ar)
		if err != nil {
			log.Println("Error marshalling track :", tr.Name)
		}
		artistJson = bytes.Trim(artistJson, "\x00")

		chanTopic := make(chan string)
		go kafka_producer.Produce("Track", trackJson, chanTopic)
		go kafka_producer.Produce("Album", albumJson, chanTopic)
		go kafka_producer.Produce("Artist", artistJson, chanTopic)
	}
}
