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
	res, err := c.client.Search().
		Index("artist_events").
		Request(&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"artist": {Query: artistName},
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
	for _, h := range hits.Hits.Hits {
		log.Println("Hits : ", h)
		artists = append(artists, h.Source)
	}

	//r := make([]ElasticDocs, len(hits.Hits.Hits))
	//for i, e := range hits.Hits.Hits {
	//	r[i].Artist = e.event.Artist
	//	r[i].Date = e.event.Date
	//	r[i].Time = e.event.Time
	//	r[i].Place = e.event.Place
	//	r[i].PlaceLink = e.event.PlaceLink
	//	r[i].Title = e.event.Title
	//	r[i].TitleLink = e.event.TitleLink
	//	r[i].Cost = e.event.Cost
	//}
	//fmt.Printf("Data with type %T from artist : %s\n", r, r)

	//var j interface{}
	//err = json.NewDecoder(res.Body).Decode(&j)
	//if err != nil {
	//	return fmt.Errorf("can't decode for artist with name : %s, error : %s", artistName, err.Error())
	//}
	//log.Println("Full j text : ", j)

	return artists, nil
}
