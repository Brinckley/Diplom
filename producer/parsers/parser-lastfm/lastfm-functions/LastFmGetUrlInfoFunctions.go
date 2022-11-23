package lastfm_functions

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

var CommonRoot = "/home/alex/docker-start/producer/"

const commonLastfmFolder = "parsers/parser-lastfm/LastFm-URL-Files/"
const commonMethods = "methods"
const commonAPI = "API"
const commonURL = "URL"

func getAPI() string {

	file, err := os.Open(CommonRoot + commonLastfmFolder + commonAPI)
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

func getMethod() []string {
	file, err := os.ReadFile(CommonRoot + commonLastfmFolder + commonMethods)
	if err != nil {
		fmt.Println("Error opening methods file")
		log.Fatal(err)
	}
	methods := strings.Split(string(file), "\n")
	return methods
}

func getURL() []string {
	file, err := os.ReadFile(CommonRoot + commonLastfmFolder + commonURL)
	if err != nil {
		fmt.Println("Error opening methods file")
		log.Fatal(err)
	}
	url := strings.Split(string(file), "\n")
	return url
}
