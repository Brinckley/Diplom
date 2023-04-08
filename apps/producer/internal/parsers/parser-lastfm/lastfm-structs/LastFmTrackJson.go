package lastfm_structs

import (
	"producer/internal/utils"
	"strconv"
)

type LastFmTrackJson struct {
	AlbumId    string
	Streamable struct {
		Fulltrack string `json:"fulltrack"`
		Text      string `json:"#text"`
	} `json:"streamable"`
	Duration *int   `json:"duration"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Attr     struct {
		Rank int `json:"rank"`
	} `json:"@attr"`
	Artist struct {
		Url  string `json:"url"`
		Name string `json:"name"`
		Mbid string `json:"mbid"`
	} `json:"artist"`
}

func NewLastFrmTrackJson(albumId string, lt *LastFmTrackJson) LastFmTrackJson {
	return LastFmTrackJson{
		AlbumId:    albumId,
		Streamable: lt.Streamable,
		Duration:   lt.Duration,
		Url:        lt.Url,
		Name:       lt.Name,
		Attr:       lt.Attr,
		Artist:     lt.Artist,
	}
}

func (l LastFmTrackJson) GetPosition() string {
	return ""
}

func (l LastFmTrackJson) GetTitle() string {
	return l.Name
}

func (l LastFmTrackJson) GetDuration() string {
	if l.Duration == nil {
		return "null"
	}
	return strconv.Itoa(*l.Duration)
}

func (l LastFmTrackJson) GetUrl() string {
	return utils.DecodeUrl(l.Url)
}

func (l LastFmTrackJson) GetArtistId() string {
	return l.Artist.Mbid
}

func (l LastFmTrackJson) GetAlbumId() string {
	return l.AlbumId
}

type OriginalLastFmTrackStruct struct {
	Track struct {
		Name       string `json:"name"`
		Mbid       string `json:"mbid"`
		URL        string `json:"url"`
		Duration   string `json:"duration"`
		Streamable struct {
			Text      string `json:"#text"`
			Fulltrack string `json:"fulltrack"`
		} `json:"streamable"`
		Listeners string `json:"listeners"`
		Playcount string `json:"playcount"`
		Artist    struct {
			Name string `json:"name"`
			Mbid string `json:"mbid"`
			URL  string `json:"url"`
		} `json:"artist"`
		Album struct {
			Artist string `json:"artist"`
			Title  string `json:"title"`
			Mbid   string `json:"mbid"`
			URL    string `json:"url"`
			Image  []struct {
				Text string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
			Attr struct {
				Position string `json:"position"`
			} `json:"@attr"`
		} `json:"album"`
		Toptags struct {
			Tag []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"tag"`
		} `json:"toptags"`
		Wiki struct {
			Published string `json:"published"`
			Summary   string `json:"summary"`
			Content   string `json:"content"`
		} `json:"wiki"`
	} `json:"track"`
}
