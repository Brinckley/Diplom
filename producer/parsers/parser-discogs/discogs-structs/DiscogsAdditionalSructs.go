package discogs_structs

import "producer/parsers/interfaces"

type DiscogsPagesReleases struct {
	Pagination struct {
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"per_page"`
		Items   int `json:"items"`
		Urls    struct {
			Last string `json:"last"`
			Next string `json:"next"`
		} `json:"urls"`
	} `json:"pagination"`
	Releases []DiscogsReleases `json:"releases"`
}

type DiscogsReleases struct {
	ID          int    `json:"id"`
	Status      string `json:"status,omitempty"`
	Type        string `json:"type"`
	Format      string `json:"format,omitempty"`
	Label       string `json:"label,omitempty"`
	Title       string `json:"title"`
	ResourceURL string `json:"resource_url"`
	Role        string `json:"role"`
	Artist      string `json:"artist"`
	Year        int    `json:"year"`
	Thumb       string `json:"thumb"`
	Stats       struct {
		Community struct {
			InWantlist   int `json:"in_wantlist"`
			InCollection int `json:"in_collection"`
		} `json:"community"`
	} `json:"stats"`
	MainRelease int `json:"main_release,omitempty"`
}

func (r DiscogsReleases) GetTracks() []interfaces.ITrack {
	var tracks []interfaces.ITrack
	for _, t := range tracks {
		tracks = append(tracks, t)
	}
	return tracks
}

func (r DiscogsReleases) GetTitle() string {
	return r.Title
}

func (r DiscogsReleases) GetTracksLen() int {
	return -1
}

func (r DiscogsReleases) GetUrl() string {
	return r.ResourceURL
}

func (r DiscogsReleases) GetYear() int {
	return r.Year
}

func (r DiscogsReleases) GetImage() string {
	return ""
}

type DiscogsMasters struct {
	ID                   int     `json:"id"`
	MainRelease          int     `json:"main_release"`
	MostRecentRelease    int     `json:"most_recent_release"`
	ResourceURL          string  `json:"resource_url"`
	URI                  string  `json:"uri"`
	VersionsURL          string  `json:"versions_url"`
	MainReleaseURL       string  `json:"main_release_url"`
	MostRecentReleaseURL string  `json:"most_recent_release_url"`
	NumForSale           int     `json:"num_for_sale"`
	LowestPrice          float64 `json:"lowest_price"`
	Images               []struct {
		Type        string `json:"type"`
		URI         string `json:"uri"`
		ResourceURL string `json:"resource_url"`
		URI150      string `json:"uri150"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
	} `json:"images"`
	Genres    []string           `json:"genres"`
	Year      int                `json:"year"`
	Tracklist []DiscogsTrackJson `json:"tracklist"`
	Artists   []struct {
		Name        string `json:"name"`
		Anv         string `json:"anv"`
		Join        string `json:"join"`
		Role        string `json:"role"`
		Tracks      string `json:"tracks"`
		ID          int    `json:"id"`
		ResourceURL string `json:"resource_url"`
	} `json:"artists"`
	Title       string `json:"title"`
	DataQuality string `json:"data_quality"`
}

func (m DiscogsMasters) GetTitle() string {
	return m.Title
}

func (m DiscogsMasters) GetTracksLen() int {
	return len(m.Tracklist)
}

func (m DiscogsMasters) GetTracks() []interfaces.ITrack {
	var tracks []interfaces.ITrack
	for _, t := range m.Tracklist {
		tracks = append(tracks, t)
	}
	return tracks
}

func (m DiscogsMasters) GetUrl() string {
	return m.URI
}

func (m DiscogsMasters) GetYear() int {
	return m.Year
}

func (m DiscogsMasters) GetImage() string {
	if len(m.Images) != 0 {
		return m.Images[0].URI
	}
	return ""
}
