package lastfm_functions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"producer/parsers/parser-lastfm/lastfm-structs"
	"strings"
)

func ReadArtist(ArtistName string) lastfm_structs.LastFmArtistJson {
	// JSON: /2.0/?method=artist.getinfo&artist=Cher&api_key=YOUR_API_KEY&format=json
	//input CommonRoot
	api := getAPI()
	methods := getMethod()
	url := getURL()
	fullUrl := url[0] + url[1] + methods[0] + url[2] + api + url[3] + ArtistName + url[6]
	fullUrl = strings.Replace(fullUrl, " ", "%20", -1)
	log.Println(fullUrl)
	response, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Something went wrong with connecting to lastfm")
		log.Fatal(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Something went wrong with reading data...")
		log.Fatal(err)
	}
	var artist lastfm_structs.LastFmArtistJson
	json.Unmarshal(data, &artist)
	//data, err = json.Marshal(&artist)
	//if err != nil {
	//	fmt.Println("Something went wrong with marshalling data...")
	//	log.Fatal(err)
	//}
	return artist
}

func ReadAlbum(ArtistName, AlbumName string) lastfm_structs.LastFmAlbumJson {
	//JSON: /2.0/?method=album.getinfo&api_key=YOUR_API_KEY&artist=Cher&album=Believe&format=json
	api := getAPI()
	methods := getMethod()
	url := getURL()
	fullUrl := url[0] + url[1] + methods[1] + url[2] + api + url[3] + ArtistName + url[4] + AlbumName + url[6]
	fullUrl = strings.Replace(fullUrl, " ", "%20", -1)
	log.Println(fullUrl)
	response, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Something went wrong with connecting to lastfm")
		log.Fatal(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Something went wrong with reading data...")
		log.Fatal(err)
	}

	var album lastfm_structs.LastFmAlbumJson
	json.Unmarshal(data, &album)

	// log.Println("NAME OF LASTFM OBJECT TEST : ", album.Album.Name)
	//data, err = json.Marshal(&album)
	//if err != nil {
	//	fmt.Println("Something went wrong with marshalling data...")
	//	log.Fatal(err)
	//}

	return album
}

func ReadTrack(ArtistName, TrackName string) lastfm_structs.LastFmTrackJson {
	// JSON: /2.0/?method=track.getInfo&api_key=YOUR_API_KEY&artist=cher&track=believe&format=json
	api := getAPI()
	methods := getMethod()
	url := getURL()
	fullUrl := url[0] + url[1] + methods[2] + url[2] + api + url[3] + ArtistName + url[5] + TrackName + url[6]
	fullUrl = strings.Replace(fullUrl, " ", "%20", -1)
	log.Println(fullUrl)
	response, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Something went wrong with connecting to lastfm")
		log.Fatal(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Something went wrong with reading data...")
		log.Fatal(err)
	}

	var track lastfm_structs.LastFmTrackJson
	json.Unmarshal(data, &track)
	//data, err = json.Marshal(&track)
	//if err != nil {
	//	fmt.Println("Something went wrong with marshalling data...")
	//	log.Fatal(err)
	//}

	return track
}
