package lastfm_structs

type LastFmArtistJson struct {
	Artist struct {
		Name  string `json:"name"`
		Mbid  string `json:"mbid"`
		Url   string `json:"url"`
		Image []struct {
			Text string `json:"#text"`
			Size string `json:"size"`
		} `json:"image"`
		Streamable string `json:"streamable"`
		Ontour     string `json:"ontour"`
		Stats      struct {
			Listeners string `json:"listeners"`
			Playcount string `json:"playcount"`
		} `json:"stats"`
		Similar struct {
			Artist []struct {
				Name  string `json:"name"`
				Url   string `json:"url"`
				Image []struct {
					Text string `json:"#text"`
					Size string `json:"size"`
				} `json:"image"`
			} `json:"artist"`
		} `json:"similar"`
		Tags struct {
			Tag []struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"tag"`
		} `json:"tags"`
		Bio struct {
			Links struct {
				Link struct {
					Text string `json:"#text"`
					Rel  string `json:"rel"`
					Href string `json:"href"`
				} `json:"link"`
			} `json:"links"`
			Published string `json:"published"`
			Summary   string `json:"summary"`
			Content   string `json:"content"`
		} `json:"bio"`
	} `json:"artist"`
}

func (a LastFmArtistJson) GetName() string {
	return a.Artist.Name
}

func (a LastFmArtistJson) GetBio() string {
	return a.Artist.Bio.Summary
}

func (a LastFmArtistJson) GetAlbumNum() int {
	return -1
}

func (a LastFmArtistJson) GetOnTour() bool {
	return a.Artist.Ontour == "1"
}

func (a LastFmArtistJson) GetImage() string {
	if len(a.Artist.Image) != 0 {
		return a.Artist.Image[0].Text
	}
	return ""
}

func (a LastFmArtistJson) GetUrl() string {
	return a.Artist.Url
}
