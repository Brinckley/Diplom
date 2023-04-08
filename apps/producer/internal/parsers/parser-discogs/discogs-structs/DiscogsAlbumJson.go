package discogs_structs

import (
	"producer/internal/parsers/interfaces"
	"producer/internal/utils"
	"strconv"
)

type DiscogsAlbumJson struct {
	ID          int    `json:"id"`
	Status      string `json:"status"`
	Year        int    `json:"year"`
	ResourceURL string `json:"resource_url"`
	URI         string `json:"uri"`
	Artists     []struct {
		Name        string `json:"name"`
		Anv         string `json:"anv"`
		Join        string `json:"join"`
		Role        string `json:"role"`
		Tracks      string `json:"tracks"`
		ID          int    `json:"id"`
		ResourceURL string `json:"resource_url"`
	} `json:"artists"`
	ArtistsSort string `json:"artists_sort"`
	Labels      []struct {
		Name           string `json:"name"`
		Catno          string `json:"catno"`
		EntityType     string `json:"entity_type"`
		EntityTypeName string `json:"entity_type_name"`
		ID             int    `json:"id"`
		ResourceURL    string `json:"resource_url"`
	} `json:"labels"`
	Series    []interface{} `json:"series"`
	Companies []struct {
		Name           string `json:"name"`
		Catno          string `json:"catno"`
		EntityType     string `json:"entity_type"`
		EntityTypeName string `json:"entity_type_name"`
		ID             int    `json:"id"`
		ResourceURL    string `json:"resource_url"`
	} `json:"companies"`
	Formats []struct {
		Name         string   `json:"name"`
		Qty          string   `json:"qty"`
		Descriptions []string `json:"descriptions"`
	} `json:"formats"`
	DataQuality string `json:"data_quality"`
	Community   struct {
		Have   int `json:"have"`
		Want   int `json:"want"`
		Rating struct {
			Count   int     `json:"count"`
			Average float64 `json:"average"`
		} `json:"rating"`
		Submitter struct {
			Username    string `json:"username"`
			ResourceURL string `json:"resource_url"`
		} `json:"submitter"`
		Contributors []struct {
			Username    string `json:"username"`
			ResourceURL string `json:"resource_url"`
		} `json:"contributors"`
		DataQuality string `json:"data_quality"`
		Status      string `json:"status"`
	} `json:"community"`
	FormatQuantity    int     `json:"format_quantity"`
	DateAdded         string  `json:"date_added"`
	DateChanged       string  `json:"date_changed"`
	NumForSale        int     `json:"num_for_sale"`
	LowestPrice       float64 `json:"lowest_price"`
	MasterID          int     `json:"master_id"`
	MasterURL         string  `json:"master_url"`
	Title             string  `json:"title"`
	Country           string  `json:"country"`
	Released          string  `json:"released"`
	Notes             string  `json:"notes"`
	ReleasedFormatted string  `json:"released_formatted"`
	Identifiers       []struct {
		Type        string `json:"type"`
		Value       string `json:"value"`
		Description string `json:"description,omitempty"`
	} `json:"identifiers"`
	Videos []struct {
		URI         string `json:"uri"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Duration    int    `json:"duration"`
		Embed       bool   `json:"embed"`
	} `json:"videos"`
	Genres       []string           `json:"genres"`
	Styles       []string           `json:"styles"`
	Tracklist    []DiscogsTrackJson `json:"tracklist"`
	Extraartists []struct {
		Name        string `json:"name"`
		Anv         string `json:"anv"`
		Join        string `json:"join"`
		Role        string `json:"role"`
		Tracks      string `json:"tracks"`
		ID          int    `json:"id"`
		ResourceURL string `json:"resource_url"`
	} `json:"extraartists"`
	Images []struct {
		Type        string `json:"type"`
		URI         string `json:"uri"`
		ResourceURL string `json:"resource_url"`
		URI150      string `json:"uri150"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
	} `json:"images"`
	Thumb           string `json:"thumb"`
	EstimatedWeight int    `json:"estimated_weight"`
	BlockedFromSale bool   `json:"blocked_from_sale"`
}

func (a *DiscogsAlbumJson) GetTitle() string {
	return a.Title
}

func (a *DiscogsAlbumJson) GetTracksLen() int {
	return len(a.Tracklist)
}

func (a *DiscogsAlbumJson) GetUrl() string {
	return utils.DecodeUrl(a.URI)
}

func (a *DiscogsAlbumJson) GetYear() int {
	return a.Year
}

func (a *DiscogsAlbumJson) GetImage() string {
	return ""
}

func (a *DiscogsAlbumJson) GetTracks() []interfaces.ITrack {
	var tracks []interfaces.ITrack
	for _, t := range a.Tracklist {
		tracks = append(tracks, t)
	}
	return tracks
}

func (a *DiscogsAlbumJson) GetArtistsId() string {
	return strconv.Itoa(a.Artists[0].ID)
}

func (a *DiscogsAlbumJson) GetAlbumId() string {
	return utils.DecodeUrl(a.ResourceURL)
}
