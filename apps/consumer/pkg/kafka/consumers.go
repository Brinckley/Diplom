package kafka

import (
	"consumer/pkg/kafka/prometheus"
	"consumer/pkg/postgres"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var DEBUG = false
var PROMETH = false

func (k *ClientKafka) ConsumeAndSend() {
	topicChan := make(chan string)
	ctx := context.Background()

	if PROMETH {
		go k.prometheusClient.StartHandling()
	}
	go k.consumeArtist(topicChan, ctx)
	go k.consumeAlbum(topicChan, ctx)
	go k.consumeTrack(topicChan, ctx)

	timeout := time.After(5000 * time.Second)
	for i := 0; i < 1; i++ {
		select {
		case tc := <-topicChan:
			log.Println("Topic Parsed :", tc)
		case <-timeout:
			log.Println("Error getting info from topic and writing to the db (timeout)!")
		}
	}
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
		//fmt.Println("[STATS] ARTIST TIME MSG : ", m.Time.Second())

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
		if !DEBUG {
			k.postgresClient.DBInsertArtist(artist) // Taken away for while debugging // adding value to postgres
		}
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
		//fmt.Println("[STATS] ALBUM TIME MSG : ", m.Time.Second())

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
		if !DEBUG {
			k.postgresClient.DBInsertAlbum(album) // Taken away for while debugging
		}
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

	timerStart := time.Now().UnixNano()
	trackCounter := 0
	var milis []int64
	for {
		trackCounter++
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
		if !DEBUG {
			k.postgresClient.DBInsertTrack(track) // Taken away for while debugging
		}

		if PROMETH {
			now := time.Now().UnixMicro()
			delta := now - timerStart
			milis = append(milis, delta)
			fmt.Println("[STATS] TRACK TIME MSG : ", delta)
			timerStart = now

			if len(milis) >= 1000 {
				log.Println("Length of milis : ", len(milis))
				for _, mil := range milis {
					k.prometheusClient.SendMessage(prometheus.NewMsgTrack(m, mil))
					time.Sleep(1 * time.Second)
				}
				timerStart = time.Now().UnixMicro()
				milis = []int64{timerStart}
			}
		}

	}
	if err := c.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Closed")

	tcs <- k.cTrackTopic
}
