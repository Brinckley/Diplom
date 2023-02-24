package discogs_functions

import (
	"io"
	"log"
	"net/http"
	"strconv"
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

func ReadArtistById(ArtistId int) []byte {
	//JSON: Get("https://api.discogs.com/artists/170355")
	fullUrl := getURL() + getMethodArtist() + "/" + strconv.Itoa(ArtistId)
	return getByteArrayByURL(fullUrl)
}

func ReadReleasesByArtistId(ArtistId int) []byte {
	fullUrl := getURL() + getMethodArtist() + "/" + strconv.Itoa(ArtistId) + "/" + getMethodReleases() + "?page=1&per_page=50000"
	log.Printf("Url for artist %v is %v\n", ArtistId, fullUrl)
	return getByteArrayByURL(fullUrl)
}

func ReadMasterById(AlbumId int) []byte {
	//JSON: Get("https://api.discogs.com/releases/249504")
	fullUrl := getURL() + getMethodMasters() + "/" + strconv.Itoa(AlbumId)
	return getByteArrayByURL(fullUrl)
}
