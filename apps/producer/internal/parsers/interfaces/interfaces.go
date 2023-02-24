package interfaces

type IArtist interface {
	GetName() string
	GetBio() string
	GetOnTour() bool
	GetImage() string
	GetUrl() string
	GetGenre() string

	GetId() string
}

type IAlbum interface {
	GetTitle() string
	GetTracksLen() int
	GetUrl() string
	GetYear() int
	GetImage() string

	GetArtistsId() string
	GetAlbumId() string

	GetTracks() []ITrack // needed for parsing
}

type ITrack interface {
	GetTitle() string
	GetDuration() string
	GetUrl() string
	GetPosition() string

	GetArtistId() string
	GetAlbumId() string
}
