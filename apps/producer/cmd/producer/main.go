package main

import (
	"producer/internal/collector"
	"producer/internal/consul"
	"producer/internal/kafka-producer"
)

func main() {
	kafka_producer.InitProducer()
	consulDebug := false

	if !consulDebug {
		artists := []string{"Skepticism", "ДДТ", "Melechesh", "Rammstein"}
		for _, a := range artists {
			collector.ParserCollectorArtistWithReleases(a)
		}
	} else {
		consul.RegisterServer()
	}

}
