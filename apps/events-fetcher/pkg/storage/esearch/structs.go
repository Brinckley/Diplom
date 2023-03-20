package esearch

import (
	"encoding/json"
	kassir_structs "events-fetcher/pkg/parsers/kassir-parser/kassir-structs"
	"log"
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
