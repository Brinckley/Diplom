package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"producer/collector"
	"producer/parsers/parser-discogs/discogs-functions"
	"producer/parsers/parser-lastfm/lastfm-functions"
	"time"

	"github.com/segmentio/kafka-go"
)

func produce(ctx context.Context, network, brokerAddress, topic string, partition int, message []byte) {
	conn, err := kafka.DialLeader(ctx, network, brokerAddress, topic, partition)
	if err != nil {
		fmt.Println("DialLeader command escaped....")
		log.Fatal(err)
	}
	conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	conn.WriteMessages(kafka.Message{Value: message})
}

//func produceFull(ctx context.Context, network, brokerAddress, topicArtist, topicAlbum, topicTrack string, partition int,
//	artist collector.ArtistDB, album collector.AlbumDB, track collector.TrackDB) {
//	conn, err := kafka.DialLeader(ctx, network, brokerAddress, topicArtist, partition)
//	if err != nil {
//		fmt.Println("DialLeader command escaped....")
//		log.Fatal(err)
//	}
//	conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
//	conn.WriteMessages(kafka.Message{Value: message})
//}

var cBrokerAddress = "localhost:9092"
var cNetwork = "tcp"
var cArtistTopic = "Artist"
var cAlbumTopic = "Album"
var cTrackTopic = "Track"

func InitProducer() {
	cBrokerAddress = os.Getenv("BROKER_ADDRESS")
	cNetwork = os.Getenv("NETWORK")
	cArtistTopic = os.Getenv("ARTIST_TOPIC_NAME")
	cAlbumTopic = os.Getenv("ALBUM_TOPIC_NAME")
	cTrackTopic = os.Getenv("TRACK_TOPIC_NAME")
	lastfm_functions.CommonRoot = os.Getenv("COMMON_ROOT")
	discogs_functions.CommonRoot = os.Getenv("COMMON_ROOT")
}

func main() {
	InitProducer()

	ArtistName := "Skepticism"
	// ArtistName
	// AlbumName
	// TrackName

	ArtistJson, AlbumsJson, TracksJson := collector.ParserCollectorArtistWithReleases(ArtistName)
	fmt.Println("....Artist With Albums And Tracks Parsed....")
	log.Println("---------------------------------------------------------------------------------------------------------------------------")

	produce(context.Background(), cNetwork, cBrokerAddress, cArtistTopic, 0, ArtistJson)
	fmt.Println("----Artist written----")
	for i, data := range AlbumsJson {
		produce(context.Background(), cNetwork, cBrokerAddress, cAlbumTopic, 0, data)
		fmt.Printf("%v.  Sent to Kafka\n", i+1)
	}
	fmt.Println("----Albums written----")
	for i, data2 := range TracksJson {
		for j, data := range data2 {
			produce(context.Background(), cNetwork, cBrokerAddress, cTrackTopic, 0, data)
			fmt.Printf("%v. %v. Sent to Kafka  (List length: %v, Album Length: %v)\n", i, j, len(TracksJson), len(data2))
		}
	}
	fmt.Println("----Tracks written----")
	log.Println("---------------------------------------------------------------------------------------------------------------------------")
	//var art collector.ArtistDB
	//json.Unmarshal(ArtistJson, &art)
	//log.Printf("Artist data: %v\n", art)
	//log.Println("---------------------------------------------------------------------------------------------------------------------------")
	//log.Println()
	//for i, data := range AlbumsJson {
	//	var alb collector.AlbumDB
	//	json.Unmarshal(data, &alb)
	//	log.Printf("%v.  %v\n", i, alb)
	//}
	//
	//var tra TrackDB
	//json.Unmarshal(TrackJson, &tra)
	//fmt.Println("TRACK \n", tra)

	//produce(context.Background(), cNetwork, cBrokerAddress, cArtistTopic, 0, ArtistJson)
	//fmt.Println("----Artist written----")
	//produce(context.Background(), cNetwork, cBrokerAddress, cAlbumTopic, 0, AlbumJson)
	//fmt.Println("----Album written-----")
	//produce(context.Background(), cNetwork, cBrokerAddress, cTrackTopic, 0, TrackJson)
	//fmt.Println("----Track written-----")
}
