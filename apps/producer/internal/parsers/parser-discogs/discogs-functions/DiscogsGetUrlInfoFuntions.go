package discogs_functions

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

var CommonRoot = "/home/alex/docker-start/producer/"

const commonDiscogsFolder = "internal/parsers/parser-discogs/Discogs-URL-Files/"
const commonMethods = "methods"
const commonAPI = "API"
const commonURL = "URL"
const commonSearch = "search"

func getAPI() string {
	file, err := os.Open(CommonRoot + commonDiscogsFolder + commonAPI)
	if err != nil {
		fmt.Println("Error opening API file")
		log.Fatal(err)
	}
	defer file.Close()

	wr := bytes.Buffer{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}
	return wr.String()
}

func getMethodArtist() string {
	file, err := os.ReadFile(CommonRoot + commonDiscogsFolder + commonMethods)
	if err != nil {
		fmt.Println("Error opening methods file")
		log.Fatal(err)
	}
	methods := strings.Split(string(file), "\n")
	return methods[0]
}

func getMethodReleases() string {
	file, err := os.ReadFile(CommonRoot + commonDiscogsFolder + commonMethods)
	if err != nil {
		fmt.Println("Error opening methods file")
		log.Fatal(err)
	}
	methods := strings.Split(string(file), "\n")
	return methods[1]
}

func getMethodMasters() string {
	file, err := os.ReadFile(CommonRoot + commonDiscogsFolder + commonMethods)
	if err != nil {
		fmt.Println("Error opening methods file")
		log.Fatal(err)
	}
	methods := strings.Split(string(file), "\n")
	return methods[2]
}

func getURL() string {
	file, err := os.ReadFile(CommonRoot + commonDiscogsFolder + commonURL)
	if err != nil {
		fmt.Println("Error opening methods file")
		log.Fatal(err)
	}
	url := strings.Split(string(file), "\n")
	return url[0]
}

func GetSearch() []string {
	file, err := os.Open(CommonRoot + commonDiscogsFolder + commonSearch)
	if err != nil {
		log.Fatal("Error reading Search info")
	}
	defer file.Close()

	var url []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url = append(url, scanner.Text())
	}
	return url
}
