package kafka

import (
	"consumer/pkg/kafka/prometheus"
	"consumer/pkg/postgres"
	"github.com/sirupsen/logrus"
	"os"
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

	exchanger chan prometheus.MsgTrack

	logger           *logrus.Logger
	postgresClient   *postgres.ClientPostgres
	prometheusClient *prometheus.ClientPrometheus
}

func NewKafka(pc *postgres.ClientPostgres, logger *logrus.Logger, prom *prometheus.ClientPrometheus, e chan prometheus.MsgTrack) *ClientKafka {
	var k ClientKafka
	k.logger = logger
	k.postgresClient = pc
	k.prometheusClient = prom
	k.exchanger = e
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
