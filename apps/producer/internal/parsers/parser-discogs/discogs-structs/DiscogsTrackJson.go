package discogs_structs

type DiscogsTrackJson struct {
	Position string `json:"position"`
	Type     string `json:"type_"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

func (t DiscogsTrackJson) GetTitle() string {
	return t.Title
}

func (t DiscogsTrackJson) GetDuration() string {
	return t.Duration
}

func (t DiscogsTrackJson) GetUrl() string {
	return ""
}

func (t DiscogsTrackJson) GetPosition() string {
	return t.Position
}

func (t DiscogsTrackJson) GetArtistId() string {
	return ""
}

func (t DiscogsTrackJson) SetAlbumId(albumId string) {
}

func (t DiscogsTrackJson) GetAlbumId() string {
	return ""
}
