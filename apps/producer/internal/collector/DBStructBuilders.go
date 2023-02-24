package collector

import (
	"hash/fnv"
	"producer/internal/parsers/interfaces"
	"strconv"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func hash(s string) uint32 { // bigint
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0
	}
	return h.Sum32()
}

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

func ArtistDBBuilder(la, da interfaces.IArtist) ArtistDB {
	artist := ArtistDB{
		Name:       da.GetName(),
		Bio:        la.GetBio(),
		OnTour:     la.GetOnTour(),
		Picture:    la.GetImage(),
		UrlLastfm:  la.GetUrl(),
		UrlDiscogs: da.GetUrl(),
		Genre:      la.GetGenre(),

		IdLastfm:  la.GetId(),
		IdDiscogs: da.GetId(),

		ArtistHash: hash(da.GetId()),
	}

	return artist
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

func AlbumDBBuilder(la, da interfaces.IAlbum) AlbumDB {
	l := max(da.GetTracksLen(), la.GetTracksLen())
	album := AlbumDB{
		Name:       da.GetTitle(),
		Release:    strconv.Itoa(da.GetYear()),
		UrlLastfm:  la.GetUrl(),
		UrlDiscogs: da.GetUrl(),
		Picture:    la.GetImage(),
		TrackCount: l,

		ArtistHash: hash(da.GetArtistsId()),
		AlbumHash:  hash(la.GetAlbumId()),
	}

	return album
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

func TrackDBBuilder(la interfaces.ITrack, da interfaces.IAlbum) TrackDB {
	track := TrackDB{
		Name:      la.GetTitle(),
		Duration:  la.GetDuration(),
		Position:  la.GetPosition(),
		UrlLastfm: la.GetUrl(),

		ArtistHash: hash(da.GetArtistsId()),
		AlbumHash:  hash(la.GetAlbumId()),
	}

	return track
}
