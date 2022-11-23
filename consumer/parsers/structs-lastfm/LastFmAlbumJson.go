package structs_lastfm

type StructAlbum struct {
	Album struct {
		Artist string `json:"artist"`
		Mbid   string `json:"mbid"`
		Tags   struct {
			Tag []struct {
				Url  string `json:"url"`
				Name string `json:"name"`
			} `json:"tag"`
		} `json:"tags"`
		Playcount string `json:"playcount"`
		Image     []struct {
			Size string `json:"size"`
			Text string `json:"#text"`
		} `json:"image"`
		Tracks struct {
			Track []struct {
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
			} `json:"track"`
		} `json:"tracks"`
		Url       string `json:"url"`
		Name      string `json:"name"`
		Listeners string `json:"listeners"`
		Wiki      struct {
			Published string `json:"published"`
			Summary   string `json:"summary"`
			Content   string `json:"content"`
		} `json:"wiki"`
	} `json:"album"`
}
