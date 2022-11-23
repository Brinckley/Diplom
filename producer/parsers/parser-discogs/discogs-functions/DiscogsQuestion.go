package discogs_functions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type DiscogsSearchStruct struct {
	Pagination struct {
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"per_page"`
		Items   int `json:"items"`
		Urls    struct {
			Last string `json:"last"`
			Next string `json:"next"`
		} `json:"urls"`
	} `json:"pagination"`
	Results []struct {
		ID       int    `json:"id"`
		Type     string `json:"type"`
		UserData struct {
			InWantlist   bool `json:"in_wantlist"`
			InCollection bool `json:"in_collection"`
		} `json:"user_data"`
		MasterID    interface{} `json:"master_id"`
		MasterURL   interface{} `json:"master_url"`
		URI         string      `json:"uri"`
		Title       string      `json:"title"`
		Thumb       string      `json:"thumb"`
		CoverImage  string      `json:"cover_image"`
		ResourceURL string      `json:"resource_url"`
		Country     string      `json:"country,omitempty"`
		Year        string      `json:"year,omitempty"`
		Format      []string    `json:"format,omitempty"`
		Label       []string    `json:"label,omitempty"`
		Genre       []string    `json:"genre,omitempty"`
		Style       []string    `json:"style,omitempty"`
		Barcode     []string    `json:"barcode,omitempty"`
		Catno       string      `json:"catno,omitempty"`
		Community   struct {
			Want int `json:"want"`
			Have int `json:"have"`
		} `json:"community,omitempty"`
		FormatQuantity int `json:"format_quantity,omitempty"`
		Formats        []struct {
			Name         string   `json:"name"`
			Qty          string   `json:"qty"`
			Text         string   `json:"text"`
			Descriptions []string `json:"descriptions"`
		} `json:"formats,omitempty"`
	} `json:"results"`
}

func CreateRequestByName(ObjName, Type string) int {
	// https://api.discogs.com/database/search?q=Nirvana&token=abcxyz123456
	searchUrl := GetSearch()
	API := getAPI()
	fullUrl := searchUrl[0] + searchUrl[1] + ObjName + searchUrl[2] + API

	var client http.Client
	resp, err := client.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var answer DiscogsSearchStruct
		json.Unmarshal(bodyBytes, &answer)

		for _, elem := range answer.Results {
			if strings.Compare(elem.Type, Type) == 0 {
				return elem.ID
			}
		}
		return -1
	}
	log.Println("Error 404")
	return -1
}
