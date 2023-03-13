package esearch

import (
	"context"
	"encoding/json"
	kassir_structs "events-fetcher/pkg/parsers/kassir-parser/kassir-structs"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
	"reflect"
)

type ElasticDocs struct {
	ID        int    `json:"id"`
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

type ESClient struct {
	address string
	client  *elasticsearch.TypedClient
}

func NewESClient() *ESClient {
	var c ESClient
	err := c.init()
	if err != nil {
		log.Fatalln("[ERR] can't connect to to elasticsearch")
		return nil
	}
	return &c
}

func (c *ESClient) init() error {
	var err error
	c.address = os.Getenv("EL_ADDRESS")
	cfg := elasticsearch.Config{Addresses: []string{c.address}}
	c.client, err = elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return err
	}
	return err
}

func (c *ESClient) AddDocument(id int, docs []kassir_structs.EventInfo) {
	log.Println("[INFO] Starting iteration over docs")
	for _id, _doc := range docs {
		_doc.ID = id + _id
		resp, err := c.client.Index("index_name").
			Request(_doc).
			Do(context.Background())
		if err != nil {
			log.Printf("[ERR] can't send info about %s to es. Error : %s\n", _doc.Artist, err.Error())
			continue
		}

		var j interface{}
		err = json.NewDecoder(resp.Body).Decode(&j)
		if err != nil {
			log.Printf("[ERR] can't decode response : %s\n", err.Error())
			continue
		}
		log.Printf("[INFO] %v. Responce body : %s\n", _id, j)

	}
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
