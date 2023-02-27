package tests

import (
	kassir_functions "events-fetcher/pkg/parsers/kassir-parser/kassir-functions"
	_ "github.com/anaskhan96/soup"
	"testing"
)

func TestLength(t *testing.T) {
	html, err := kassir_functions.LookAtSearchHtml("ддт", "рок")
	if err != nil {
		return
	}

	wantLength := 1
	if len(html.EventInfo) != wantLength {
		t.Fatalf("Want : %v, but got : %v", wantLength, len(html.EventInfo))
	}
}

func TestTitle(t *testing.T) {
	html, err := kassir_functions.LookAtSearchHtml("ддт", "рок")
	if err != nil {
		return
	}

	wantTitle := "ДДТ"
	if html.EventInfo[0].Title != wantTitle {
		t.Fatalf("Want : %v, but got : %v", wantTitle, html.EventInfo[0].Title)
	}
}

func TestTitleLink(t *testing.T) {
	html, err := kassir_functions.LookAtSearchHtml("ддт", "рок")
	if err != nil {
		return
	}

	wantTitleLink := "https://msk.kassir.ru/koncert/ddt"
	if html.EventInfo[0].TitleLink != wantTitleLink {
		t.Fatalf("Want : %v, but got : %v", wantTitleLink, html.EventInfo[0].TitleLink)
	}
}

func TestDate(t *testing.T) {
	html, err := kassir_functions.LookAtSearchHtml("ддт", "рок")
	if err != nil {
		return
	}

	wantDate := "17 Июнь"
	if html.EventInfo[0].Date != wantDate {
		t.Fatalf("Want : %v, but got : %v", wantDate, html.EventInfo[0].Date)
	}
}

func TestTime(t *testing.T) {
	html, err := kassir_functions.LookAtSearchHtml("ддт", "рок")
	if err != nil {
		return
	}

	wantTime := "сб 20:00"
	if html.EventInfo[0].Time != wantTime {
		t.Fatalf("Want : %v, but got : %v", wantTime, html.EventInfo[0].Time)
	}
}

func TestPlace(t *testing.T) {
	html, err := kassir_functions.LookAtSearchHtml("ддт", "рок")
	if err != nil {
		return
	}

	wantPlace := "Стадион «Открытие Банк Арена»"
	if html.EventInfo[0].Place != wantPlace {
		t.Fatalf("Want : %v, but got : %v", wantPlace, html.EventInfo[0].Place)
	}
}

func TestPlaceLink(t *testing.T) {
	html, err := kassir_functions.LookAtSearchHtml("ддт", "рок")
	if err != nil {
		return
	}

	wantPlaceLink := "https://msk.kassir.ru/sportivnye-kompleksy/stadion-otkryitie-arena-5001"
	if html.EventInfo[0].PlaceLink != wantPlaceLink {
		t.Fatalf("Want : %v, but got : %v", wantPlaceLink, html.EventInfo[0].PlaceLink)
	}
}

func TestCost(t *testing.T) {
	html, err := kassir_functions.LookAtSearchHtml("ддт", "рок")
	if err != nil {
		return
	}

	wantCost := "2 900 — 20 000"
	if html.EventInfo[0].Cost != wantCost {
		t.Fatalf("Want : %v, but got : %v", wantCost, html.EventInfo[0].Cost)
	}
}
