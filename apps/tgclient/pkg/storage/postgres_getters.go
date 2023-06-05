package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

func (p *TgPostgres) GetAlbumsByArtist(artistName string) ([]string, error) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		return []string{}, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	queryGetByName := fmt.Sprintf("select %s.* from %s join %s "+
		"on %s.artistHash=%s.artistHash and %s.name='%s' order by albums.id ASC",
		p.cTableAlbum, p.cTableAlbum, p.cTableArtists,
		p.cTableAlbum, p.cTableArtists, p.cTableArtists,
		artistName)

	log.Println(queryGetByName)

	rows, err := conn.Query(context.Background(), queryGetByName)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var albumList []string
	innerId := 1
	for rows.Next() {
		var album AlbumDB
		err := rows.Scan(&album.Id, &album.Name, &album.Release, &album.UrlLastfm, &album.UrlDiscogs,
			&album.Picture, &album.TrackCount, &album.ArtistHash, &album.AlbumHash)
		if err != nil {
			fmt.Println(err)
			continue
		}

		album.Id = innerId
		albumList = append(albumList, album.ToString())
	}

	return albumList, nil
}
