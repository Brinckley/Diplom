package discogs_structs

import "strconv"

type DiscogsArtistJson struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	ResourceURL string `json:"resource_url"`
	URI         string `json:"uri"`
	ReleasesURL string `json:"releases_url"`
	Images      []struct {
		Type        string `json:"type"`
		URI         string `json:"uri"`
		ResourceURL string `json:"resource_url"`
		URI150      string `json:"uri150"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
	} `json:"images"`
	Profile        string   `json:"profile"`
	Urls           []string `json:"urls"`
	Namevariations []string `json:"namevariations"`
	Aliases        []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ResourceURL string `json:"resource_url"`
	} `json:"aliases"`
	Members []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ResourceURL string `json:"resource_url"`
		Active      bool   `json:"active"`
	} `json:"members"`
	DataQuality string `json:"data_quality"`
}

func (d DiscogsArtistJson) GetName() string {
	return d.Name
}

func (d DiscogsArtistJson) GetBio() string {
	return d.Profile
}

func (d DiscogsArtistJson) GetAlbumNum() int {
	return -1
}

func (d DiscogsArtistJson) GetOnTour() bool {
	return false
}

func (d DiscogsArtistJson) GetImage() string {
	if len(d.Images) != 0 {
		return d.Images[0].URI
	}
	return ""
}

func (d DiscogsArtistJson) GetUrl() string {
	return d.URI
}

func (d DiscogsArtistJson) GetId() string {
	return strconv.Itoa(d.ID)
}

func (d DiscogsArtistJson) GetGenre() string {
	return ""
}
