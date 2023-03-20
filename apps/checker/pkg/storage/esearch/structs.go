package esearch

import (
	"encoding/json"
	"log"
	"reflect"
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

type SpecialHits struct {
	Hits struct {
		Hits []struct {
			Source ElasticDocs `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// A function for marshaling structs to JSON string
func jsonStruct(doc *ElasticDocs) (string, error) {
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
