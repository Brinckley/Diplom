package structs_lastfm

type StructArtist struct {
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
