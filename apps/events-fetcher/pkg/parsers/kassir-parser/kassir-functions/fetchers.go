package kassir_functions

import (
	"github.com/anaskhan96/soup"
	"log"
	"strings"
)

func fetchTitle(root soup.Root) string {
	path := root.
		Find("div", "class", "title").
		Find("a")
	if path.Error != nil {
		log.Println("no title data found")
		return ""
	}

	title := strings.TrimSpace(path.Text())
	return title
}

func fetchTitleLink(root soup.Root) string {
	path := root.
		Find("div", "class", "title").
		Find("a")
	if path.Error != nil {
		log.Println("no title link data found")
		return ""
	}

	titleLink := strings.TrimSpace(path.Attrs()["href"])
	return titleLink
}

func fetchDate(root soup.Root) string {
	//path := root

	dt := ""
	for _, tr := range root.FindAll("span") {
		dt += tr.Text() + " "
	}

	return dt
	//day := path.Find("span", "class", "day-m").HTML()
	//mouth := path.Find("span", "class", "mouth").HTML()
	//time := path.Find("span", "class", "time").HTML()
	//
	//if path.Error != nil {
	//	log.Println("no date data found")
	//	return ""
	//}

	//date := strings.TrimSpace(path.Text())
	//date := day + " " + mouth + " " + time
	//return date
}

func fetchTime(root soup.Root) string {
	path := root.
		Find("time", "class", "date").
		Find("span")
	if path.Error != nil {
		log.Println("no time data found")
		return ""
	}

	time := strings.TrimSpace(path.Text())
	return time
}

func fetchPlace(root soup.Root) string {
	path := root.
		Find("div", "class", "venue").
		Find("a")
	if path.Error != nil {
		log.Println("no place data found")
		return ""
	}

	place := strings.TrimSpace(path.Text())
	return place
}

func fetchPlaceLink(root soup.Root) string {
	path := root.
		Find("div", "class", "venue").
		Find("a")
	if path.Error != nil {
		log.Println("no place link data found")
		return ""
	}

	placeLink := strings.TrimSpace(path.Attrs()["href"])
	return placeLink
}

func fetchCost(root soup.Root) string {
	path := root.
		Find("div", "class", "cost").
		Find("div", "class", "price")
	if path.Error != nil {
		log.Println("no cost data found")
		return ""
	}

	cost := strings.TrimSpace(path.Text())
	return cost
}
