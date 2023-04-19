package kassir_functions

import (
	"errors"
	kstructs "events-fetcher/pkg/parsers/kassir-parser/kassir-structs"
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
	"strings"
	"time"
)

func LookAtSearchHtml(artistNameRaw, genre string) ([]kstructs.EventInfo, error) {
	artistName := strings.Replace(artistNameRaw, " ", "%20", -1) // removing spaces
	category, err := kstructs.SelectGenre(genre)
	if err != nil {
		return []kstructs.EventInfo{}, err
	}

	genreCategory := "&" + category // adding category parameter to the url
	fullUrl := fmt.Sprintf("https://msk.kassir.ru/category?main=3000%s&sort=1&c=90&keyword=%s", genreCategory, artistName)
	//log.Println("Full Search Url : ", fullUrl)

	doc, err := getHTMLFromLink(fullUrl)
	if err != nil {
		log.Println("no such url")
		return []kstructs.EventInfo{}, err
	}

	docc := doc.
		Find("div", "class", "tiles-container")
	if docc.Error != nil {
		log.Printf("[ERR] no info found about artist '%s' with genre '%s'\n", artistNameRaw, genre)
		return []kstructs.EventInfo{}, err
	}

	commonPath := docc.FindAll("div", "class", "new--w-12")
	if len(commonPath) == 0 { // no concerts at all option
		return []kstructs.EventInfo{}, errors.New(fmt.Sprintf("[ERR] no info found about artist '%s' with genre '%s'\n", artistNameRaw, genre))
	}

	//log.Printf("[INFO] starting parsing events for artist '%s'\n", artistNameRaw)
	var events []kstructs.EventInfo

	for _, p := range commonPath {
		//log.Println("----------------------------------------------111 !1111-----------", p.HTML())
		innerPath := p.
			Find("div", "class", "tile-card")
		if innerPath.Error != nil {
			log.Println("No inner info found!")
			continue
		}
		ei := kstructs.EventInfo{ // constructing event object from found data using direct html searchers
			Artist:    artistNameRaw,
			Title:     fetchTitle(innerPath),
			TitleLink: fetchTitleLink(innerPath),
			Date:      fetchDate(innerPath),
			Time:      fetchTime(innerPath),
			Place:     fetchPlace(innerPath),
			PlaceLink: fetchPlaceLink(innerPath),
			Cost:      fetchCost(innerPath),
			TimeStamp: time.Now().Unix(),
		}
		//log.Println("EVENT CONTENT : ", ei)
		events = append(events, ei)
	}
	//log.Printf("[INFO] end of parsing events for artist '%s'\n", artistNameRaw)
	return events, nil
}

func getHTMLFromLink(fullUrl string) (soup.Root, error) {
	resp, err := soup.Get(fullUrl)
	if err != nil {
		return soup.Root{}, err
	}
	doc := soup.HTMLParse(resp)
	return doc, nil // returning all html code from the page
}
