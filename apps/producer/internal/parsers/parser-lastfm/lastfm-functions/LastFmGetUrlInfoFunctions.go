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

const commonLastfmFolder = "internal/parsers/parser-lastfm/LastFm-URL-Files/"
const commonMethods = "methods"
const commonAPI = "API"
const commonURL = "URL"

func readStringsArray(path string) []string {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error opening file :", path)
		log.Fatal(err)
	}
	array := strings.Split(string(file), "\n")
	return array
}

func getApi() string {
	path := CommonRoot + commonLastfmFolder + commonAPI
	file, err := os.Open(path)
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

func getMethodArtistGetInfo() string {
	path := CommonRoot + commonLastfmFolder + commonMethods
	methods := readStringsArray(path)
	return methods[0]
}
func getMethodAlbumGetInfo() string {
	path := CommonRoot + commonLastfmFolder + commonMethods
	methods := readStringsArray(path)
	return methods[1]
}
func getMethodTrackGetInfo() string {
	path := CommonRoot + commonLastfmFolder + commonMethods
	methods := readStringsArray(path)
	return methods[2]
}

func getURLMain() string {
	path := CommonRoot + commonLastfmFolder + commonURL
	url := readStringsArray(path)
	return url[0]
}
func getURLMethod() string {
	path := CommonRoot + commonLastfmFolder + commonURL
	url := readStringsArray(path)
	return url[1]
}
func getURLApi() string {
	path := CommonRoot + commonLastfmFolder + commonURL
	url := readStringsArray(path)
	return url[2]
}
func getURLArtist() string {
	path := CommonRoot + commonLastfmFolder + commonURL
	url := readStringsArray(path)
	return url[3]
}
func getURLAlbum() string {
	path := CommonRoot + commonLastfmFolder + commonURL
	url := readStringsArray(path)
	return url[4]
}
func getURLTrack() string {
	path := CommonRoot + commonLastfmFolder + commonURL
	url := readStringsArray(path)
	return url[5]
}
func getURLFormat() string {
	path := CommonRoot + commonLastfmFolder + commonURL
	url := readStringsArray(path)
	return url[6]
}
