package lastfm_structs

import "strconv"

type LastFmTrackJson struct {
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

func (l LastFmTrackJson) GetPosition() string {
	return ""
}

func (l LastFmTrackJson) GetTitle() string {
	return l.Name
}

func (l LastFmTrackJson) GetDuration() string {
	return strconv.Itoa(*l.Duration)
}

func (l LastFmTrackJson) GetUrl() string {
	return l.Url
}
