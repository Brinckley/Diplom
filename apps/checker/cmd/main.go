package main

import (
	"checker/pkg/kafka"
	"checker/pkg/storage/esearch"
	"checker/pkg/storage/postgres"
	"context"
	"log"
	"sync"
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

	esclient := esearch.NewESClient() // connect to esearch, where we will get info about every event
	ekafka := kafka.NewKafka()

	eventsChan := make(chan []esearch.ElasticDocs, 1)
	var mutex sync.Mutex
	cnt := len(names)
	wg := sync.WaitGroup{}

	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func(name string) {
			log.Println("\nSearching for name : ", name)
			events, err := esclient.SearchArtist(name)
			if err != nil {
				log.Fatalln("[ERR] ", err.Error())
			}
			eventsChan <- events
		}(names[i])
	}

	for i := 0; i < cnt; i++ {
		go ReceiveFormChan(eventsChan, ekafka, &wg, &mutex)
	}
	wg.Wait()
}

func ReceiveFormChan(c chan []esearch.ElasticDocs, ekafka *kafka.ClientKafka, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer func() {
		mutex.Unlock()
		wg.Done()
	}()
	mutex.Lock()
	err := ekafka.ProduceEvents(context.Background(), <-c)
	if err != nil {
		log.Println("[ERR] can't send the message to kafka : ", err.Error())
	}
}
