package esearch

import (
	"bytes"
	"context"
	"encoding/json"
	kassir_structs "events-fetcher/pkg/parsers/kassir-parser/kassir-structs"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
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
		log.Fatalln("[ERR] can't connect to to elasticsearch") // no need to work if es is absent
		return nil
	}
	return &c
}

func (c *ESClient) init() error {
	var err error
	c.address = os.Getenv("EL_ADDRESS") // getting address from docker-compose
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

func (c *ESClient) AddDocument(docs []kassir_structs.EventInfo) {
	if len(docs) > 0 {
		log.Println("[ESEARCH INFO] artist name for sending : ", docs[0].Artist)
	} else {
		log.Println("[ESEARCH INFO] empty event list")
		return
	}
	//log.Println("[ESEARCH INFO] starting iteration over docs")
	//log.Println("[ESEARCH INFO] artist name : ", docs[0].Artist)
	for _, _doc := range docs {
		res, err := c.client.Search().
			Size(1000).
			Index("artist_events").
			Request(&search.Request{
				Query: &types.Query{
					Match: map[string]types.MatchQuery{
						"artist": {Query: _doc.Artist},
					},
				},
			}).Do(context.Background())
		if res != nil {
			log.Printf("[ESEARCH INFO] match found for artist : '%s', Event : '%s', Hits : '%v'\n",
				_doc.Artist, _doc.Title, len(res.Hits.Hits))
			exists := false
			for _, hit := range res.Hits.Hits {
				var ei kassir_structs.EventInfo
				err := json.Unmarshal(hit.Source_, &ei)
				if err != nil {
					log.Println("[ERR] can't unmarshall hit : ", hit.Source_)
					continue
				}
				if ei.TitleLink == _doc.TitleLink {
					exists = true
					break
				}
			}
			if exists {
				continue
			} else {
				log.Printf("[ESEARCH INFO] match NOT found for artist : '%s', Event : '%s'\n", _doc.Artist, _doc.Title)
			}
		} else {
			log.Printf("[ESEARCH INFO] match NOT found for artist : '%s', Event : '%s'\n", _doc.Artist, _doc.Title)
		}

		docJson, err := json.Marshal(_doc)
		if err != nil {
			log.Printf("[ERR] can't marshall event : '%v', err : '%s'\f", _doc, err.Error())
			continue
		}

		req := esapi.IndexRequest{
			Index:   "artist_events",
			Body:    bytes.NewReader(docJson),
			Refresh: "true",
		}
		// Perform the request with the client.
		resReq, err := req.Do(context.Background(), c.client)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		func() { _ = resReq.Body.Close() }()

		if err != nil {
			log.Printf("[ERR] can't send info about %s to es. Error : %s\n", _doc.Artist, err.Error())
			continue
		}

		log.Printf("[ESEARCH INFO] response after sending data of artist : '%s'\nResponse status : '%s'\nHas warnings param : '%v'\n",
			_doc.Artist, resReq.Status(), resReq.HasWarnings())
	}
}
