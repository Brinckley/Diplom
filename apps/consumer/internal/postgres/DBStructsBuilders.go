package postgres

type ArtistDB struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Bio        string `json:"bio"`
	OnTour     bool   `json:"onTour"`
	Picture    string `json:"picture"`
	UrlLastfm  string `json:"urlLastfm"`
	UrlDiscogs string `json:"urlDiscogs"`
	Genre      string `json:"genre"`

	IdLastfm  string `json:"idLastfm"`
	IdDiscogs string `json:"idDiscogs"`

	ArtistHash uint32 `json:"artistHash"`
}

type AlbumDB struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Release    string `json:"release"`
	UrlLastfm  string `json:"urlLastfm"`
	UrlDiscogs string `json:"urlDiscogs"`
	Picture    string `json:"picture"`
	TrackCount int    `json:"trackCount"`

	ArtistHash uint32 `json:"artistHash"`
	AlbumHash  uint32 `json:"albumHash"`
}

type TrackDB struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	UrlLastfm string `json:"urlLastfm"`
	Duration  string `json:"duration"`
	Position  string `json:"position"`

	ArtistHash uint32 `json:"artistHash"`
	AlbumHash  uint32 `json:"albumHash"`
}
