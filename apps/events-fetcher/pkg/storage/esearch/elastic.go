package esearch

import (
	"context"
	"encoding/json"
	kassir_structs "events-fetcher/pkg/parsers/kassir-parser/kassir-structs"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ElasticDocs struct {
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	TitleLink string `json:"titleLink"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Place     string `json:"place"`
	PlaceLink string `json:"placeLink"`
	Cost      string `json:"cost"`
}

func EventToElastic(event kassir_structs.EventInfo) ElasticDocs {
	return ElasticDocs{
		Artist:    event.Artist,
		Title:     event.Title,
		TitleLink: event.TitleLink,
		Date:      event.Date,
		Time:      event.Time,
		Place:     event.Place,
		PlaceLink: event.Place,
		Cost:      event.PlaceLink,
	}
}

var (
	address string
	docMap  map[string]interface{} // Create a mapping for the Elasticsearch documents
)

func NewESClient() (*elasticsearch.Client, error) {
	address = os.Getenv("EL_ADDRESS")
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}
	log.Printf("[INFO] address is : %s", address)
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("[ERR] elasticsearch connection error: %s\n", err)
		return nil, err
	}
	// Have the client instance return a response
	res, err := client.Info()
	if err != nil {
		log.Printf("client.Info() ERROR: %s", err)
		return nil, err
	} else {
		log.Printf("client response: %s", res)
	}

	if err != nil {
		return nil, err
	}
	log.Println("Initialized successfully")
	log.SetFlags(0) // Allow for custom formatting of log output
	log.Println("docMap:", docMap)
	log.Println("docMap TYPE:", reflect.TypeOf(docMap))

	return client, err
}

func AddDocument(ctx context.Context, client *elasticsearch.Client, id int, docs []kassir_structs.EventInfo) error {
	log.Println("[INFO] Starting iteration over docs")
	for _id, _doc := range docs {
		doc := EventToElastic(_doc)
		log.Println("Got doc : ", doc)
		bod, err := jsonStruct(doc)
		if err != nil {
			log.Println("Something went wrong when converting to json")
			return err
		}

		// Instantiate a request object
		req := esapi.IndexRequest{
			Index:      "some_index",
			DocumentID: strconv.Itoa(_id + id + 1),
			Body:       strings.NewReader(bod),
			Refresh:    "true",
		}
		fmt.Println("Type of request : ", reflect.TypeOf(req))

		// Return an API response object from request
		log.Println("Doing request for id :", _id+id+1)
		log.Println("[Content] : ", doc)
		log.Printf("[CONTENT] Type of client : %T, client value : %v\n", client, &client)
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer func() { _ = res.Body.Close() }()

		log.Println("[STEP] Before error check")
		log.Printf("[STEP] type of result %T\n", res)
		log.Printf("[STEP] value of result %v\n", res)

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), _id+id+1)
		} else {
			// Deserialize the response into a map.
			log.Println("[STEP] In else block (error clear)")
			var resMap map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				log.Printf("\nIndexRequest() RESPONSE:")
				// Print the response status and indexed document version.
				log.Println("Status:", res.Status())
				log.Println("Result:", resMap["result"])
				log.Println("Version:", int(resMap["_version"].(float64)))
				log.Println("resMap:", resMap)
				log.Println()
			}
		}
	}

	log.Println("[STEP] Returning nil")
	return nil
}

// A function for marshaling structs to JSON string
func jsonStruct(doc ElasticDocs) (string, error) {
	// Create struct instance of the Elasticsearch fields struct object
	docStruct := &ElasticDocs{
		Artist:    doc.Artist,
		Title:     doc.Title,
		TitleLink: doc.TitleLink,
		Date:      doc.Date,
		Time:      doc.Time,
		Place:     doc.Place,
		PlaceLink: doc.PlaceLink,
		Cost:      doc.Cost,
	}

	log.Println("\ndocStruct:", docStruct)
	log.Println("docStruct TYPE:", reflect.TypeOf(docStruct))

	// Marshal the struct to JSON and check for errors
	b, err := json.Marshal(docStruct)
	if err != nil {
		log.Println("json.Marshal ERROR:", err)
		return "", err
	}

	return string(b), nil
}
