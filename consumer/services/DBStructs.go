package services

type ArtistDB struct {
	Id         int    `json:"id"`
	AName      string `json:"aName"`
	Bio        string `json:"bio"`
	AlbumNum   int    `json:"albumNum"`
	OnTour     bool   `json:"onTour"`
	Picture    string `json:"picture"`
	UrlLastfm  string `json:"urlLastfm"`
	UrlDiscogs string `json:"urlDiscogs"`

	Albums []int `json:"albums"`
	Tracks []int `json:"tracks"`
}

type AlbumDB struct {
	Id         int    `json:"id"`
	AName      string `json:"aName"`
	Release    string `json:"release"`
	UrlLastfm  string `json:"urlLastfm"`
	UrlDicogs  string `json:"urlDicogs"`
	Picture    string `json:"picture"`
	TrackCount int    `json:"trackCount"`

	Artists []int `json:"artists"`
	Tracks  []int `json:"tracks"`
}

type TrackDB struct {
	Id        int    `json:"id"`
	TName     string `json:"tName"`
	UrlLastfm string `json:"urlLastfm"`
	Duration  string `json:"duration"`

	Artists []int `json:"artists"`
	Albums  []int `json:"albums"`
}
