package collector

import (
	"producer/parsers/interfaces"
	"strconv"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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

func ArtistDBBuilder(la, da interfaces.IArtist, n int) ArtistDB {
	artist := ArtistDB{
		AName:      da.GetName(),
		Bio:        la.GetBio(),
		AlbumNum:   n,
		OnTour:     la.GetOnTour(),
		Picture:    la.GetImage(),
		UrlLastfm:  la.GetUrl(),
		UrlDiscogs: da.GetUrl(),
	}

	return artist
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

func AlbumDBBuilder(la, da interfaces.IAlbum) AlbumDB {
	l := max(da.GetTracksLen(), la.GetTracksLen())
	album := AlbumDB{
		AName:      da.GetTitle(),
		Release:    strconv.Itoa(da.GetYear()),
		UrlLastfm:  la.GetUrl(),
		UrlDicogs:  da.GetUrl(),
		Picture:    la.GetImage(),
		TrackCount: l,
	}

	return album
}

type TrackDB struct {
	Id        int    `json:"id"`
	TName     string `json:"tName"`
	UrlLastfm string `json:"urlLastfm"`
	Duration  string `json:"duration"`

	Artists []int `json:"artists"`
	Albums  []int `json:"albums"`
}

func TrackDBBuilder(la interfaces.ITrack) TrackDB {
	track := TrackDB{
		TName:     la.GetTitle(),
		Duration:  la.GetDuration(),
		UrlLastfm: la.GetUrl(),
	}

	return track
}
