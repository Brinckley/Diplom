package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var cDialect = ""
var cHost = ""
var cDBPort = ""
var cDBUser = ""
var cDBPassword = ""
var cDBName = ""
var cTableNameArtist = ""
var cTableNameAlbum = ""
var cTableNameTrack = ""

var err error

func InitDatabase() {
	cDBPort = os.Getenv("DB_PORT")
	cDBPassword = os.Getenv("DB_PASSWORD")
	cDBName = os.Getenv("DB_NAME")
	cTableNameArtist = os.Getenv("DB_NAME_ARTIST")
	cTableNameAlbum = os.Getenv("DB_NAME_ALBUM")
	cTableNameTrack = os.Getenv("DB_NAME_TRACK")
}

func PGXInsert(dsnL, dsnR string) {
	dsn := dsnL + cDBPassword + dsnR
	artist, album, track := ReadAll(context.Background())
	log.Println("Kafka data read...")
	log.Println("===================================================================================================")

	log.Println("\nStarting db connection...")
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	log.Println("Connected to database...")

	var id int

	err = conn.QueryRow(context.Background(),
		"insert into artists(aname, bio, albumNum, onTour, picture, urlLastfm, urlDiscogs) values ($1, $2, $3, $4, $5, $6, $7) returning id;",
		artist.AName, artist.Bio, artist.AlbumNum, artist.OnTour, artist.Picture, artist.UrlLastfm, artist.UrlDicogs).Scan(&id)
	CheckError(err, "Artists")
	fmt.Println("Artist inserted id : ", id)
	err = conn.QueryRow(context.Background(),
		"insert into albums(aname, release, picture, trackCount, urlLastfm, urlDiscogs) values ($1, $2, $3, $4, $5, $6) returning id;",
		album.AName, album.Release, album.Picture, album.TrackCount, album.UrlDicogs, album.UrlLastfm).Scan(&id)
	CheckError(err, "Albums")
	fmt.Println("Album inserted id :", id)
	err = conn.QueryRow(context.Background(),
		"insert into tracks(tname, release, lyrics, urlLastfm) values ($1, $2, $3, $4) returning id;",
		track.TName, track.Release, track.Lyrics, track.UrlLastfm).Scan(&id)
	CheckError(err, "Tracks")
	fmt.Println("Track inserted id :", id)

	log.Println("....End of queries....")
	log.Println("===================================================================================================")
}

func CheckError(err error, db string) {
	if err != nil {
		fmt.Println("Failed connecting to table :", db)
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to the table :", db)
	}
}
