package kassir_functions

import (
	"errors"
	structs2 "events-fetcher/internal/parsers/kassir-parser/kassir-structs"
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
)

func LookAtSearchHtml(artistName, genre string) (structs2.Events, error) {
	category, err := structs2.SelectGenre(genre)
	if err != nil {
		return structs2.Events{}, err
	}

	genreCategory := "&" + category
	fullUrl := fmt.Sprintf("https://msk.kassir.ru/category?main=3000%s&sort=1&c=90&keyword=%s", genreCategory, artistName)
	log.Println("Full Search Url : ", fullUrl)

	doc, err := getHTMLFromLink(fullUrl)
	if err != nil {
		log.Println("no such url")
		return structs2.Events{}, err
	}

	docc := doc.
		Find("div", "class", "tiles-container")
	if docc.Error != nil {
		log.Println("no info found")
		return structs2.Events{}, err
	}

	commonPath := docc.FindAll("div", "class", "new--w-12")
	if len(commonPath) == 0 {
		return structs2.Events{}, errors.New("no info found")
	}

	var events []structs2.EventInfo
	for _, p := range commonPath {
		innerPath := p.
			Find("div", "class", "event-card__caption")
		if innerPath.Error != nil {
			log.Println("No inner info found!")
			continue
		}
		ei := structs2.EventInfo{
			Artist:    artistName,
			Title:     fetchTitle(innerPath),
			TitleLink: fetchTitleLink(innerPath),
			Date:      fetchDate(innerPath),
			Time:      fetchTime(innerPath),
			Place:     fetchPlace(innerPath),
			PlaceLink: fetchPlaceLink(innerPath),
			Cost:      fetchCost(innerPath),
		}
		events = append(events, ei)
	}
	return structs2.Events{EventInfo: events}, nil
}

func getHTMLFromLink(fullUrl string) (soup.Root, error) {
	resp, err := soup.Get(fullUrl)
	if err != nil {
		return soup.Root{}, err
	}
	doc := soup.HTMLParse(resp)
	return doc, nil
}
