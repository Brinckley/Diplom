package esearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
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
	var v ESClient
	err := v.init()
	if err != nil {
		log.Fatalf("[ERR] can't initialize elasticsearch client!")
		return nil
	}
	return &v
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

func (c *ESClient) SearchArtist(artistName string) ([]ElasticDocs, error) {
	log.Println("-----------------Searching for name : ", artistName)
	res, err := c.client.Search().
		Size(1000).
		Index("artist_events").
		Request(&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"title": {Query: artistName},
				},
			},
		}).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't search for artist with name : %s, error : %s", artistName, err.Error())
	}

	log.Println("Content : ", res)
	log.Println("Header : ", res.Header)
	log.Println("Body : ", res.Body)

	var hits SpecialHits
	err = json.NewDecoder(res.Body).Decode(&hits) // XXX: error omitted
	if err != nil {
		return nil, fmt.Errorf("can't decode for artist with name : %s, error : %s", artistName, err.Error())
	}

	var artists []ElasticDocs
	for i, h := range hits.Hits.Hits {
		log.Printf("Hit %v : %v\n", i, h)
		artists = append(artists, h.Source)
	}

	return artists, nil
}
