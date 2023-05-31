package postgres

import (
	"github.com/sirupsen/logrus"
	"os"
)

type ClientPostgres struct {
	cDBPassword      string
	cTableNameArtist string
	cTableNameAlbum  string
	cTableNameTrack  string
	cDsnL            string
	cDsnR            string
	cDsn             string

	logger *logrus.Logger
}

var err error

func NewPostgres(logger *logrus.Logger) *ClientPostgres {
	var cp ClientPostgres
	cp.logger = logger
	cp.init()
	return &cp
}

func (p *ClientPostgres) init() {
	p.cDBPassword = os.Getenv("DB_PASSWORD")
	p.cTableNameArtist = os.Getenv("DB_NAME_ARTIST")
	p.cTableNameAlbum = os.Getenv("DB_NAME_ALBUM")
	p.cTableNameTrack = os.Getenv("DB_NAME_TRACK")
	p.cDsnL = os.Getenv("DSN_LEFT")
	p.cDsnR = os.Getenv("DSN_RIGHT")
	p.cDsn = p.cDsnL + p.cDBPassword + p.cDsnR
}
