package interfaces

type IArtist interface {
	GetName() string
	GetBio() string
	GetAlbumNum() int
	GetOnTour() bool
	GetImage() string
	GetUrl() string
}

type IAlbum interface {
	GetTitle() string
	GetTracksLen() int
	GetUrl() string
	GetYear() int
	GetImage() string
	GetTracks() []ITrack
}

type ITrack interface {
	GetTitle() string
	GetDuration() string
	GetUrl() string
	GetPosition() string
}
