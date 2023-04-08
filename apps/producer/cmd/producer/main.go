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
		artists := []string{"Мельница", "Филипп+Киркоров", "Би-2", "Rammstein", "Nazareth"}
		for _, a := range artists {
			collector.ParserCollectorArtistWithReleases(a)
		}
	} else {
		consul.RegisterServer()
	}

}
