package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
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
var cTableNameArtistsAlbums = ""
var cTableNameAlbumsTracks = ""
var cDsnL = ""
var cDsnR = ""

var err error

func InitDatabase() {
	cDBPort = os.Getenv("DB_PORT")
	cDBPassword = os.Getenv("DB_PASSWORD")
	cDBName = os.Getenv("DB_NAME")
	cTableNameArtist = os.Getenv("DB_NAME_ARTIST")
	cTableNameAlbum = os.Getenv("DB_NAME_ALBUM")
	cTableNameTrack = os.Getenv("DB_NAME_TRACK")
	cTableNameArtistsAlbums = os.Getenv("DB_NAME_ARTIST_ALBUM")
	cTableNameAlbumsTracks = os.Getenv("DB_NAME_ALBUM_TRACK")
	cDsnL = os.Getenv("DSN_LEFT")
	cDsnR = os.Getenv("DSN_RIGHT")
}

func CheckError(err error, db string) {
	if err != nil {
		fmt.Println("Failed connecting to table :", db)
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to the table :", db)
	}
}

/*
CREATE TABLE IF NOT EXISTS artists(
      id          SERIAL   PRIMARY KEY,
      name        VARCHAR(100) NOT NULL,
      onTour      BOOLEAN,
      picture     VARCHAR(255),
      idLastfm    VARCHAR(255),
      urlLastfm   VARCHAR(255),
      idDiscogs   VARCHAR(255),
      urlDiscogs  VARCHAR(255),
      genre       VARCHAR(31),

      artistHash  INTEGER
    );
*/

func DBSelectArtists() {
	dsn := cDsnL + cDBPassword + cDsnR
	//log.Println("DSN : ", dsn)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), "select * from artists")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	log.Println("------------------VALUES FROM TABLE ARTISTS------------------")
	for rows.Next() {
		a := ArtistDB{}
		err := rows.Scan(&a.Id, &a.Name, &a.Bio, &a.OnTour, &a.Picture, &a.UrlLastfm, &a.UrlDiscogs, &a.Genre,
			&a.IdLastfm, &a.IdDiscogs, &a.ArtistHash)
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Printf("%v. Name: %v\n", a.Id, a.Name)
	}
}

func DBInsertArtist(artist ArtistDB) {
	dsn := cDsnL + cDBPassword + cDsnR
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var artistId int
	queryInsert := fmt.Sprintf("insert into %s (name, bio, onTour, picture, idLastfm, idDiscogs, genre, urlLastfm, urlDiscogs, artistHash) "+
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id;", cTableNameArtist)
	err = conn.QueryRow(context.Background(),
		queryInsert,
		artist.Name, artist.Bio, artist.OnTour, artist.Picture, artist.IdLastfm, artist.IdDiscogs, artist.Genre,
		artist.UrlLastfm, artist.UrlDiscogs, artist.ArtistHash).Scan(&artistId)
	CheckError(err, "Artists")
	log.Println("Artist with id added :", artistId)
}

/*
CREATE TABLE IF NOT EXISTS albums(
      id          SERIAL   PRIMARY KEY,
      name        VARCHAR(100) NOT NULL,
      release     VARCHAR(100),
	  urlLastfm   VARCHAR(255),
	  urlDiscogs  VARCHAR(255),
      picture     VARCHAR(255),
      trackCount  INTEGER,

      artistHash  INTEGER,
      albumHash   INTEGER
    );
*/

func DBSelectAlbums() {
	dsn := cDsnL + cDBPassword + cDsnR
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), "select * from "+cTableNameAlbum)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	log.Println("------------------VALUES FROM TABLE ALBUMS------------------")
	for rows.Next() {
		a := AlbumDB{}
		err := rows.Scan(&a.Id, &a.Name, &a.Release, &a.UrlLastfm, &a.UrlDiscogs, &a.Picture, &a.TrackCount, &a.ArtistHash, &a.AlbumHash)
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Printf("%v. Name: %v\n", a.Id, a.Name)
	}
}

func DBInsertAlbum(album AlbumDB) {
	dsn := cDsnL + cDBPassword + cDsnR
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var albumId int
	queryInsert := fmt.Sprintf("insert into %s (name, release, urlLastfm, urlDiscogs, picture, trackCount, artistHash, albumHash) "+
		"values ($1, $2, $3, $4, $5, $6, $7, $8) returning id;", cTableNameAlbum)
	err = conn.QueryRow(context.Background(),
		queryInsert,
		album.Name, album.Release, album.UrlLastfm, album.UrlDiscogs, album.Picture,
		album.TrackCount, album.ArtistHash, album.AlbumHash).Scan(&albumId)
	CheckError(err, "Albums")
	log.Println("Album with id added :", albumId)
}

/*
CREATE TABLE IF NOT EXISTS tracks(
      id          SERIAL   PRIMARY KEY,
	  name        VARCHAR(100) NOT NULL,
      urlLastfm   VARCHAR(255),
      duration    VARCHAR(100),
      position    VARCHAR(100),

      artistHash  INTEGER,
      albumHash   INTEGER
    );
*/

func DBSelectTracks() {
	dsn := cDsnL + cDBPassword + cDsnR
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), "select * from "+cTableNameTrack)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	log.Println("------------------VALUES FROM TABLE TRACKS------------------")
	for rows.Next() {
		t := TrackDB{}
		err := rows.Scan(&t.Id, &t.Name, &t.UrlLastfm, &t.Duration, &t.Position, &t.ArtistHash, &t.AlbumHash)
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Printf("%v. Name: %v\n", t.Id, t.Name)
	}
}

func DBInsertTrack(track TrackDB) {
	dsn := cDsnL + cDBPassword + cDsnR
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var trackId int
	queryInsert := fmt.Sprintf("insert into %s (name, urlLastfm, duration, position, artistHash, albumHash) "+
		"values ($1, $2, $3, $4, $5, $6) returning id;", cTableNameTrack)
	err = conn.QueryRow(context.Background(),
		queryInsert,
		track.Name, track.UrlLastfm, track.Duration, track.Position, track.ArtistHash, track.AlbumHash).Scan(&trackId)
	CheckError(err, "Tracks")
	log.Println("Track with id added :", trackId)
}

/*
CREATE TABLE artists_albums
(
    id        serial                                             not null unique,
    artist_id int references artists (id) on delete cascade      not null,
    album_id  int references albums (id) on delete cascade       not null
);


CREATE TABLE albums_tracks
(
    id       serial                                              not null unique,
    album_id int references albums (id) on delete cascade        not null,
    track_id int references tracks (id) on delete cascade        not null
);
*/
