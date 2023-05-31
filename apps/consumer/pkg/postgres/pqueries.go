package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

func CheckError(err error, db string) {
	if err != nil {
		fmt.Println("Failed connecting to table :", db)
		log.Fatal(err)
	} else {
		//fmt.Println("Successfully connected to the table :", db)
	}
}

func (p *ClientPostgres) DBSelectArtists() {
	//log.Println("DSN : ", dsn)
	conn, err := pgx.Connect(context.Background(), p.cDsn)
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

func (p *ClientPostgres) DBInsertArtist(artist ArtistDB) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		log.Fatalf("[ERR] unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	var artistId int
	queryInsert := fmt.Sprintf("insert into %s (name, bio, onTour, picture, idLastfm, idDiscogs, genre, urlLastfm, urlDiscogs, artistHash) "+
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id;", p.cTableNameArtist)
	err = conn.QueryRow(context.Background(),
		queryInsert,
		artist.Name, artist.Bio, artist.OnTour, artist.Picture, artist.IdLastfm, artist.IdDiscogs, artist.Genre,
		artist.UrlLastfm, artist.UrlDiscogs, artist.ArtistHash).Scan(&artistId)
	CheckError(err, "Artists")
	//log.Println("[POSTGRESQL] artist with id added :", artistId)
}

func (p *ClientPostgres) DBSelectAlbums() {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		log.Fatalf("[ERR] unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	rows, err := conn.Query(context.Background(), "select * from "+p.cTableNameAlbum)
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

func (p *ClientPostgres) DBInsertAlbum(album AlbumDB) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		log.Fatalf("[ERR] unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	var albumId int
	queryInsert := fmt.Sprintf("insert into %s (name, release, urlLastfm, urlDiscogs, picture, trackCount, artistHash, albumHash) "+
		"values ($1, $2, $3, $4, $5, $6, $7, $8) returning id;", p.cTableNameAlbum)
	err = conn.QueryRow(context.Background(),
		queryInsert,
		album.Name, album.Release, album.UrlLastfm, album.UrlDiscogs, album.Picture,
		album.TrackCount, album.ArtistHash, album.AlbumHash).Scan(&albumId)
	CheckError(err, "Albums")
	//log.Println("[POSTGRESQL] album with id added :", albumId)
}

func (p *ClientPostgres) DBSelectTracks() {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		log.Fatalf("[ERR] unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	rows, err := conn.Query(context.Background(), "select * from "+p.cTableNameTrack)
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

func (p *ClientPostgres) DBInsertTrack(track TrackDB) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		log.Fatalf("[ERR] unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	var trackId int
	queryInsert := fmt.Sprintf("insert into %s (name, urlLastfm, duration, position, artistHash, albumHash) "+
		"values ($1, $2, $3, $4, $5, $6) returning id;", p.cTableNameTrack)
	err = conn.QueryRow(context.Background(),
		queryInsert,
		track.Name, track.UrlLastfm, track.Duration, track.Position, track.ArtistHash, track.AlbumHash).Scan(&trackId)
	CheckError(err, "Tracks")
	//log.Println("[POSTGRESQL] track with id added :", trackId)
}
