package esearch

import (
	"context"
	"encoding/json"
	kassir_structs "events-fetcher/pkg/parsers/kassir-parser/kassir-structs"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

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

func (c *ESClient) CheckConnection() {
	log.Println("[INFO] Check elasticsearch current condition : ", c.client.Info())
}

func (c *ESClient) AddDocument(id int, docs []kassir_structs.EventInfo) {
	if len(docs) > 0 {
		log.Println("Artist name for sending : ", docs[0].Artist)
	}
	log.Println("[INFO] Starting iteration over docs")
	for _id, _doc := range docs {
		_doc.ID = id + _id
		resp, err := c.client.Index("artist_events").
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
		log.Printf("[INFO] %v. Responce body : %s\n", _id+1, j)
	}
}
