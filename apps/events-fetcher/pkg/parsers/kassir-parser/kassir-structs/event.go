package kassir_structs

import "encoding/json"

type EventInfo struct {
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	TitleLink string `json:"titleLink"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Place     string `json:"place"`
	PlaceLink string `json:"placeLink"`
	Cost      string `json:"cost"`
}

type Events struct {
	EventInfo []EventInfo `json:"eventInfo"`
}

func (e *Events) GetJson() ([]byte, error) {
	eventsJson, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return eventsJson, nil
}

func (e *Events) GetArtist() string {
	if len(e.EventInfo) != 0 {
		return e.EventInfo[0].Artist
	}
	return ""
}
