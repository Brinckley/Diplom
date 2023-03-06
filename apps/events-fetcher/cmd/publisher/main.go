package main

import (
	"context"
	kassir_functions "events-fetcher/pkg/parsers/kassir-parser/kassir-functions"
	kassir_structs "events-fetcher/pkg/parsers/kassir-parser/kassir-structs"
	"events-fetcher/pkg/storage/esearch"
	"events-fetcher/pkg/storage/postgres"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"net/http"
	"sync"
)

func main() {
	// reading ONE/ALL??? artist from postgres
	// sending ONE artist to parse ALL related events from kassir.ru
	// sending ALL related events to elastic

	_, err := http.NewRequest("GET", "http://elasticsearch:9200", nil)

	if err != nil {
		log.Println("!!!!!!!!!!!! Err check : ", err.Error())
	} else {
		log.Println("CHECK PASSED!!!!!!!!")
	}

	p := postgres.NewPostgres()
	p.Init()
	names, genres, err := p.GetArtistsNamesGenres()
	if err != nil {
		return
	}
	for i := 0; i < len(names); i++ {
		fmt.Printf("%s  -  %s\n", names[i], genres[i])
	}
	esclient, err := esearch.NewESClient()
	if err != nil {
		log.Fatalf("[ERR] can't connect to elastic : %s", err.Error())
	}
	log.Printf("[CONTENT] Type of client : %T, client value : %v\n", esclient, esclient)

	log.Println("------------------- All artist from db have been read -------------------")

	eventsChan := make(chan []kassir_structs.EventInfo, 1)
	var mutex sync.Mutex
	cnt := len(names)
	wg := sync.WaitGroup{}

	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func(name string, genre string) {
			html, err := kassir_functions.LookAtSearchHtml(name, genre)
			if err != nil {
				wg.Done()
				return
			}
			eventsChan <- html
		}(names[i], genres[i])
	}
	for i := 0; i < cnt; i++ {
		go ReceiveFormChan(eventsChan, i, &wg, &mutex, esclient)
	}
	wg.Wait()
}

func ReceiveFormChan(c chan []kassir_structs.EventInfo, id int, wg *sync.WaitGroup, mutex *sync.Mutex, client *elasticsearch.Client) {
	defer func() {
		mutex.Unlock()
		wg.Done()
	}()
	mutex.Lock()
	err := esearch.AddDocument(context.Background(), client, id, <-c)
	if err != nil {
		return
	}
}

// еще один топик для выхода
// чеккер прозванивает, через кафку отправляет
