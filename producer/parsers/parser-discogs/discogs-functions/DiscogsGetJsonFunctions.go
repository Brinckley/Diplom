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

func ReadArtistById(ArtistId int) []byte {
	//JSON: Get("https://api.discogs.com/artists/170355")
	methods := getMethod()
	url := getURL()
	fullUrl := url[0] + methods[0] + "/" + strconv.Itoa(ArtistId)

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

func ReadReleasesByArtistId(ArtistId int) []byte {
	methods := getMethod()
	url := getURL()
	fullUrl := url[0] + methods[0] + "/" + strconv.Itoa(ArtistId) + "/" + methods[1] + "?page=1&per_page=5000"
	log.Printf("Url for artist %v is %v\n", ArtistId, fullUrl)

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
		return bodyBytes
	}
	// log.Println("Error 404")
	return []byte{}
}

func ReadMasterById(AlbumId int) []byte {
	//JSON: Get("https://api.discogs.com/releases/249504")
	methods := getMethod()
	url := getURL()
	fullUrl := url[0] + methods[2] + "/" + strconv.Itoa(AlbumId)

	resp, err := Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(fullUrl)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//bodyString := string(bodyBytes)
		return bodyBytes
	}
	log.Println("Error 404")
	return []byte{}
}
