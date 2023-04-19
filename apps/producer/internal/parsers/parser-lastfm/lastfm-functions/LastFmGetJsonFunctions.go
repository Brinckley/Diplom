package lastfm_functions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	lastfm_structs2 "producer/internal/parsers/parser-lastfm/lastfm-structs"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func Get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return Client.Do(request)
}

func getByteArrayByURL(fullUrl string) []byte {
	resp, err := Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//bodyString := string(bodyBytes)
		return bodyBytes
	}
	//log.Println("Error 404")
	return []byte{}
}

func ReadArtist(ArtistName string) lastfm_structs2.LastFmArtistJson {
	// JSON: /2.0/?method=artist.getinfo&artist=Cher&api_key=YOUR_API_KEY&format=json
	//input CommonRoot
	ArtistName = strings.ReplaceAll(ArtistName, " ", "+")
	fullUrl := getURLMain() + getURLMethod() + getMethodArtistGetInfo() +
		getURLApi() + getApi() + getURLArtist() + ArtistName + getURLFormat()
	fullUrl = strings.Replace(fullUrl, " ", "+", -1)
	log.Printf("[INFO] lastfm full artist url : '%s'\n", fullUrl)
	data := getByteArrayByURL(fullUrl)
	var artist lastfm_structs2.LastFmArtistJson
	err := json.Unmarshal(data, &artist)
	if err != nil {
		return lastfm_structs2.LastFmArtistJson{}
	}
	return artist
}

func ReadAlbum(ArtistName, AlbumName string) lastfm_structs2.LastFmAlbumJson {
	//JSON: /2.0/?method=album.getinfo&api_key=YOUR_API_KEY&artist=Cher&album=Believe&format=json
	AlbumName = strings.ReplaceAll(AlbumName, " ", "+")
	fullUrl := getURLMain() + getURLMethod() + getMethodAlbumGetInfo() +
		getURLApi() + getApi() + getURLArtist() + ArtistName + getURLAlbum() + AlbumName + getURLFormat()
	fullUrl = strings.Replace(fullUrl, " ", "+", -1)
	//log.Println(fullUrl)
	data := getByteArrayByURL(fullUrl)
	var artist lastfm_structs2.LastFmArtistJson
	err := json.Unmarshal(data, &artist)
	if err != nil {
		return lastfm_structs2.LastFmAlbumJson{}
	}
	var album lastfm_structs2.LastFmAlbumJson
	err = json.Unmarshal(data, &album)
	if err != nil {
		return lastfm_structs2.LastFmAlbumJson{}
	}
	return album
}

func ReadTrack(ArtistName, TrackName string) lastfm_structs2.LastFmTrackJson {
	// JSON: /2.0/?method=track.getInfo&api_key=YOUR_API_KEY&artist=cher&track=believe&format=json
	ArtistName = strings.ReplaceAll(ArtistName, ",", "+")
	TrackName = strings.ReplaceAll(TrackName, ",", "+")
	fullUrl := getURLMain() + getURLMethod() + getMethodTrackGetInfo() +
		getURLApi() + getApi() + getURLArtist() + ArtistName + getURLTrack() + TrackName + getURLFormat()
	fullUrl = strings.Replace(fullUrl, " ", "+", -1)
	//log.Println(fullUrl)
	data := getByteArrayByURL(fullUrl)
	var artist lastfm_structs2.LastFmArtistJson
	json.Unmarshal(data, &artist)
	var track lastfm_structs2.LastFmTrackJson
	json.Unmarshal(data, &track)

	return track
}
