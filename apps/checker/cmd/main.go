package main

import (
	"checker/pkg/kafka"
	"checker/pkg/storage/esearch"
	"checker/pkg/storage/postgres"
	"context"
	"log"
)

func main() {
	pg := postgres.NewPostgres()       // connect to db, where we will get info about artists
	names, err := pg.GetArtistsNames() // getting lists of names
	if err != nil {
		log.Fatalf("[ERR] error getting artists' names : %s", err.Error())
		return
	}
	log.Println("[INFO] names read")
	log.Println("Names : ", names)

	ekafka := kafka.NewKafka()
	esclient := esearch.NewESClient() // connect to esearch, where we will get info about every event
	for i := 0; i < len(names); i++ {
		log.Println("\nSearching for name : ", names[i])
		events, err := esclient.SearchArtist(names[i])
		if err != nil {
			if err != nil {
				log.Fatalln("[ERR] ", err.Error())
			}
			continue
		}
		err = ekafka.ProduceEvents(context.Background(), events)
		if err != nil {
			log.Fatalln("[ERR] ", err.Error())
		}

	}
}
