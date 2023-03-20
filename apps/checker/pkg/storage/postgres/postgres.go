package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type ArtistPostgres struct {
	cDBPort       string
	cDBPassword   string
	cDBName       string
	cTableArtists string
	cDsnL         string
	cDsnR         string
}

func NewPostgres() *ArtistPostgres {
	var db ArtistPostgres
	db.init()
	return &db
}

func (p *ArtistPostgres) init() {
	p.cDBPort = os.Getenv("DB_PORT")
	p.cDBPassword = os.Getenv("DB_PASSWORD")
	p.cDBName = os.Getenv("DB_NAME")
	p.cTableArtists = os.Getenv("DB_NAME_ARTIST")
	p.cDsnL = os.Getenv("DSN_LEFT")
	p.cDsnR = os.Getenv("DSN_RIGHT")
}

func (p *ArtistPostgres) GetArtistsNames() ([]string, error) {
	dsn := p.cDsnL + p.cDBPassword + p.cDsnR
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	artistNames, err := p.getFields(conn, p.cTableArtists, "name")
	if err != nil {
		return nil, err
	}

	return artistNames, nil
}

func (p *ArtistPostgres) getFields(conn *pgx.Conn, table string, field string) ([]string, error) {
	queryGet := fmt.Sprintf("select %s from %s", field, table)
	rows, err := conn.Query(context.Background(), queryGet)
	fmt.Println("Get all : ", queryGet)
	if err != nil {
		return []string{}, err
	}
	var vals []string
	for rows.Next() {
		var val string
		err := rows.Scan(&val)
		if err != nil {
			fmt.Println(err)
			continue
		}
		vals = append(vals, val)
	}

	return vals, nil
}
