package main

import (
	kassir_functions "events-fetcher/internal/parsers/kassir-parser/kassir-functions"
	"fmt"
)

func main() {
	abbr, err := kassir_functions.CityAbbr("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Abbr :", abbr)
}
