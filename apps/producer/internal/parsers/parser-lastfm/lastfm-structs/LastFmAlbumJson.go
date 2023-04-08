package lastfm_structs

import (
	"producer/internal/parsers/interfaces"
	"producer/internal/utils"
)

type LastFmAlbumJson struct {
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
			Track []LastFmTrackJson `json:"track"`
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

func (a LastFmAlbumJson) GetTitle() string {
	return a.Album.Name
}

func (a LastFmAlbumJson) GetTracksLen() int {
	return len(a.Album.Tracks.Track)
}

func (a LastFmAlbumJson) GetUrl() string {
	return utils.DecodeUrl(a.Album.Url)
}

func (a LastFmAlbumJson) GetYear() int {
	return 0
}

func (a LastFmAlbumJson) GetImage() string {
	if len(a.Album.Image) != 0 {
		return a.Album.Image[0].Text
	}
	return ""
}

func (a LastFmAlbumJson) GetTracks() []interfaces.ITrack {
	var tracks []interfaces.ITrack
	for _, t := range a.Album.Tracks.Track {
		tCopy := NewLastFrmTrackJson(a.GetAlbumId(), &t)
		tracks = append(tracks, tCopy)
	}
	return tracks
}

func (a LastFmAlbumJson) GetArtistsId() string {
	s := ""
	s = s + ", " + a.Album.Artist
	return s
}

func (a LastFmAlbumJson) GetAlbumId() string {
	return a.Album.Mbid
}
