package main

import (
	kassir_functions "events-fetcher/pkg/parsers/kassir-parser/kassir-functions"
	kassir_structs "events-fetcher/pkg/parsers/kassir-parser/kassir-structs"
	"events-fetcher/pkg/storage/esearch"
	"events-fetcher/pkg/storage/postgres"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	// reading ONE/ALL??? artist from postgres
	// sending ONE artist to parse ALL related events from kassir.ru
	// sending ALL related events to elastic

	log.Println("----------------------------Temporary check----------------------------")
	_, err := http.NewRequest("GET", "http://elasticsearch:9200", nil)
	if err != nil {
		log.Fatalln("!!!!!!!!!!!! Err check : ", err.Error())
	} else {
		log.Println("CHECK PASSED!!!!!!!!")
	}
	log.Println("----------------------------End of temporary check----------------------------")

	p := postgres.NewPostgres()
	p.Init()
	names, genres, err := p.GetArtistsNamesGenres() // writing all artists' names into list
	if err != nil {
		log.Fatalln("[ERR] can't get data about artists from postgres : ", err.Error())
	}
	log.Println("----------------------------All artists with their genres----------------------------")
	for i := 0; i < len(names); i++ {
		fmt.Printf("%s  -  %s\n", names[i], genres[i])
	}
	log.Println("----------------------------End of all artists----------------------------")

	esclient := esearch.NewESClient()
	if err != nil {
		log.Fatalf("[ERR] can't connect to elastic : %s", err.Error())
	}
	log.Printf("[CONTENT] Type of client : %T, client value : %v\n", esclient, esclient)
	esclient.CheckConnection()

	log.Println("------------------- All artist from db have been read, es connection set -------------------")

	eventsChan := make(chan []kassir_structs.EventInfo, 1)
	var mutex sync.Mutex
	cnt := len(names)
	wg := sync.WaitGroup{}
	counter := 0

	for i := 0; i < cnt; i++ {
		wg.Add(1)                            // adding work for each of artists
		go func(name string, genre string) { // func to get all concerts of one artist in selected genre
			allEvents, err := kassir_functions.LookAtSearchHtml(name, genre)
			if err != nil {
				wg.Done() // nothing to do if error exists, job done
				return
			}
			counter += len(allEvents)
			eventsChan <- allEvents // sending got events to processing for es
		}(names[i], genres[i]) // function works for each artist from the list
	}
	for i := 0; i < cnt; i++ { // for each artist we have a separate channel
		go ReceiveFormChan(eventsChan, &wg, &mutex, esclient)
	}
	wg.Wait()
	log.Println("[INFO] All EVENTS AMOUNT : ", counter)
}

func ReceiveFormChan(c chan []kassir_structs.EventInfo, wg *sync.WaitGroup, mutex *sync.Mutex, client *esearch.ESClient) {
	defer func() { // ending mutex usage to unlock es, finishing the work
		mutex.Unlock()
		wg.Done()
	}()
	mutex.Lock()
	// sending all events related to one artist to the es for adding
	client.AddDocument(<-c) // id is an inner field that basically is not needed, but let it be

}
