package main

import (
	"events-fetcher/pkg/storage/postgres"
	"fmt"
)

func main() {
	p := postgres.NewPostgres()
	p.Init()
	names, genres, err := p.GetArtistsNamesGenres()
	if err != nil {
		return
	}
	for i := 0; i < len(names); i++ {
		fmt.Printf("%s  -  %s\n", names[i], genres[i])
	}
}
