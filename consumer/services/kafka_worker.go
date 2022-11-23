package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
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

func consumeArtist(ctx context.Context, network, brokerAddress, topic string, partition int) ArtistDB {
	conn, err := kafka.DialLeader(ctx, network, brokerAddress, topic, partition)
	if err != nil {
		fmt.Println("DialLeader command escaped....")
		log.Fatal(err)
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	// var artists []json_structs.StructArtist
	batch := conn.ReadBatch(1e3, 1e9)
	byteb := make([]byte, 1e6)
	message := ""
	for {
		_, err := batch.Read(byteb)
		if err != nil {
			break
		}
		message += string(bytes.Trim(byteb, "\x00"))
	}
	// fmt.Println(message)

	var art ArtistDB
	err = json.Unmarshal([]byte(message), &art)
	if err != nil {
		log.Fatalln(err)
	}

	return art
}

func consumeAlbum(ctx context.Context, network, brokerAddress, topic string, partition int) AlbumDB {
	conn, err := kafka.DialLeader(ctx, network, brokerAddress, topic, partition)
	if err != nil {
		fmt.Println("DialLeader command escaped....")
		log.Fatal(err)
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	// var artists []json_structs.StructArtist
	batch := conn.ReadBatch(1e3, 1e9)
	byteb := make([]byte, 1e6)
	message := ""
	for {
		_, err := batch.Read(byteb)
		if err != nil {
			break
		}
		message += string(bytes.Trim(byteb, "\x00"))
	}
	// fmt.Println(message)

	var alb AlbumDB
	err = json.Unmarshal([]byte(message), &alb)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("Struct :", alb)
	return alb
}

func consumeTrack(ctx context.Context, network, brokerAddress, topic string, partition int) TrackDB {
	conn, err := kafka.DialLeader(ctx, network, brokerAddress, topic, partition)
	if err != nil {
		fmt.Println("DialLeader command escaped....")
		log.Fatal(err)
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	// var artists []json_structs.StructArtist
	batch := conn.ReadBatch(1e3, 1e9)
	byteb := make([]byte, 1e6)
	message := ""
	for {
		_, err := batch.Read(byteb)
		if err != nil {
			break
		}
		message += string(bytes.Trim(byteb, "\x00"))
	}
	// fmt.Println(message)

	var tra TrackDB
	err = json.Unmarshal([]byte(message), &tra)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("Struct :", tra)
	return tra
}

func ReadAll(context context.Context) (artist ArtistDB, album AlbumDB, track TrackDB) {
	artist = consumeArtist(context, cNetwork, cBrokerAddress, cArtistTopic, 0)
	fmt.Println("Artist read....", artist)
	album = consumeAlbum(context, cNetwork, cBrokerAddress, cAlbumTopic, 0)
	fmt.Println("Album  read....", album)
	track = consumeTrack(context, cNetwork, cBrokerAddress, cTrackTopic, 0)
	fmt.Println("Track  read....", track)

	return
}
