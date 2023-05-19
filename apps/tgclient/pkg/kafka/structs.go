package kafka

import "fmt"

type Event struct {
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	TitleLink string `json:"titleLink"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Place     string `json:"place"`
	PlaceLink string `json:"placeLink"`
	Cost      string `json:"cost"`
	TimeStamp int64  `json:"timeStamp"`
}

func (e *Event) CreateNotification() string {
	e.fillTheBlank()
	msg := fmt.Sprintf("Hi! There is a new event of your favourite artist %s.\n"+
		"It is called '%s'.\n"+
		"Date is %s, %s.\n"+
		"Place where the event will be held is %s.\n."+
		"Tickets price is %s.\n"+
		"Event link : '%s'.\n"+
		"Place link : '%s'.\n"+
		"Good Luck!",
		e.Artist, e.Title, e.Date, e.Time, e.Place, e.Cost, e.TitleLink, e.PlaceLink)
	return msg
}

func (e *Event) fillTheBlank() {
	if e.Date == "" {
		e.Date = "Date yet unknown. Check the links please for accurate information."
	}
	if e.Cost == "" {
		e.Cost = "Costs are yet unknown. Check the links please for accurate information."
	}
}
