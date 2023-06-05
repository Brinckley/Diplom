package kafka

import (
	"fmt"
)

type Event struct {
	Artist       string `json:"artist"`
	Title        string `json:"title"`
	TitleLink    string `json:"titleLink"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	Place        string `json:"place"`
	PlaceLink    string `json:"placeLink"`
	Cost         string `json:"cost"`
	TimeStamp    int64  `json:"timeStamp"`
	TimeReceived int64
	Number       int
}

func (e *Event) CreateNotification() string {
	e.fillTheBlank()
	msg := fmt.Sprintf("Привет! У твоего любимого артиста '%s' скоро концерт.\n"+
		"Название : '%s'.\n"+
		"Дата и время : %s\n"+
		"Место проведения : %s.\n"+
		"Цены : %s.\n"+
		"Ссылка на событие : '%s'.\n"+
		"Ссылка на место : '%s'.\n",
		e.Artist, e.Title, e.Date, e.Place, e.Cost, e.TitleLink, e.PlaceLink)
	return msg
}

func (e *Event) fillTheBlank() {
	if e.Date == "" {
		//e.Date = "Date yet unknown. Check the links please for accurate information."
		e.Date = " уточните дату на сайте"
	}
	if e.Cost == "" {
		//e.Cost = "Costs are yet unknown. Check the links please for accurate information."
		e.Cost = " уточните цены на сайте."
	}
}
