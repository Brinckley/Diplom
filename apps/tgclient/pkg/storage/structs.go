package storage

import (
	"fmt"
	"tgclient/pkg/utils"
)

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

func (a *AlbumDB) ToString() string {
	a.UrlLastfm = utils.DecodeUrl(a.UrlLastfm)
	a.UrlDiscogs = utils.DecodeUrl(a.UrlDiscogs)

	return fmt.Sprintf(
		"%v) %s, %s. %v треков.",
		//	"Discogs : %s\n"+
		//	"LastFm : %s\n",
		a.Id, a.Name, a.Release, a.TrackCount,
		//a.UrlDiscogs,
		//a.UrlLastfm
	)
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
