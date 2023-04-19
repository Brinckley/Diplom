package storage

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"tgclient/pkg/kafka"
)

var ErrUserExists = errors.New("user already exists")
var ErrNoArtists = errors.New("no artists found")
var ErrNoFavorites = errors.New("no favorites found")
var ErrNoSuchUser = errors.New("no such user found")
var ErrNoSuchArtist = errors.New("no such artist found")
var ErrNoSuchPair = errors.New("no such pair found")

type TgPostgres struct {
	cDBPort          string
	cDBPassword      string
	cDBName          string
	cTableArtists    string
	cTableAlbum      string
	cTableTrack      string
	cTableUsers      string
	cTableUserArtist string
	cDsnL            string
	cDsnR            string
	cDsn             string
}

func NewPostgres() *TgPostgres {
	var db TgPostgres
	db.Init()
	return &db
}

func (p *TgPostgres) Init() {
	p.cDBPort = os.Getenv("DB_PORT")
	p.cDBPassword = os.Getenv("DB_PASSWORD")
	p.cDBName = os.Getenv("DB_NAME")
	p.cTableArtists = os.Getenv("DB_NAME_ARTIST")
	p.cTableAlbum = os.Getenv("DB_NAME_ALBUM")
	p.cTableTrack = os.Getenv("DB_NAME_TRACK")
	p.cTableUsers = os.Getenv("DB_NAME_USERS")
	p.cDsnL = os.Getenv("DSN_LEFT")
	p.cDsnR = os.Getenv("DSN_RIGHT")
	p.cTableUserArtist = os.Getenv("DB_NAME_USER_ARTIST")
	p.cDsn = p.cDsnL + p.cDBPassword + p.cDsnR
}

func CheckError(err error, table string) {
	if err != nil {
		log.Println("Failed connecting to table :", table)
		log.Fatal(err)
	} else {
		log.Println("Successfully connected to the table :", table)
	}
}

func (p *TgPostgres) Registration(message *tgbotapi.Message) error {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	existence, err := checkExistence(conn, p.cTableUsers, "username", message.From.UserName)
	if err != nil {
		return err
	}
	if existence {
		return ErrUserExists
	}

	queryInsert := fmt.Sprintf("insert into %s (username, firstName, сhatID) "+
		"values ($1, $2, $3) returning id;", p.cTableUsers)
	var userId int
	err = conn.QueryRow(context.Background(),
		queryInsert,
		message.From.UserName, message.From.FirstName, message.Chat.ID).Scan(&userId)
	CheckError(err, "users")
	log.Println("User with id added :", userId)

	return nil
}

func (p *TgPostgres) GetAllArtists() ([]string, error) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	rows, err := conn.Query(context.Background(), "select name from "+p.cTableArtists)
	fmt.Println("select name from ", p.cTableArtists)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []string
	for rows.Next() {
		var artist string
		err := rows.Scan(&artist)
		if err != nil {
			fmt.Println(err)
			continue
		}
		artists = append(artists, artist)
	}

	if len(artists) == 0 {
		return nil, ErrNoArtists
	}
	return artists, nil
}

func (p *TgPostgres) GetFavorites(message *tgbotapi.Message) ([]string, error) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	querySelect := fmt.Sprintf(""+
		"select %s.name from %s inner join %s on %s.id=%s.artist_id "+
		"inner join %s on %s.id = %s.user_id and %s.username='%s';",
		p.cTableArtists, p.cTableArtists, p.cTableUserArtist, p.cTableArtists, p.cTableUserArtist,
		p.cTableUsers, p.cTableUsers, p.cTableUserArtist, p.cTableUsers, message.From.UserName)
	rows, err := conn.Query(context.Background(), querySelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favs []string
	for rows.Next() {
		var fav string
		err := rows.Scan(&fav)
		if err != nil {
			return nil, err
		}
		favs = append(favs, fav)
	}

	if len(favs) == 0 {
		return nil, ErrNoFavorites
	}
	return favs, nil
}

func (p *TgPostgres) GetAllSubscribers(event kafka.Event) ([]int64, error) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	querySelectSubs := fmt.Sprintf(
		"select %s.сhatid from %s "+
			"join %s on %s.id = %s.user_id "+
			"join %s on %s.id = %s.artist_id and %s.name = '%s' "+
			"and %s.updTime < %v",
		p.cTableUsers, p.cTableUsers,
		p.cTableUserArtist, p.cTableUsers, p.cTableUserArtist,
		p.cTableArtists, p.cTableArtists, p.cTableUserArtist, p.cTableArtists,
		event.Artist, p.cTableUserArtist, event.TimeStamp)
	log.Println("Get All Subscriber query : ", querySelectSubs)

	rows, err := conn.Query(context.Background(), querySelectSubs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []int64
	for rows.Next() {
		var sub int64
		err := rows.Scan(&sub)
		if err != nil {
			return nil, err
		}
		subscribers = append(subscribers, sub)
	}

	err = p.updateSubscriptionTimeStamp(subscribers)
	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

func (p *TgPostgres) getId(conn *pgx.Conn, table string, field string, value string) (int, error) {
	queryGetId := fmt.Sprintf("select id from %s where %s='%s'", table, field, value)
	fmt.Println("Get Id : ", queryGetId)
	rows, err := conn.Query(context.Background(), queryGetId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var res int
	for rows.Next() {
		err := rows.Scan(&res)
		if err != nil {
			return -1, err
		}
	}
	return res, nil
}

func checkExistence(conn *pgx.Conn, table string, field string, value string) (bool, error) {
	query := fmt.Sprintf("select EXISTS (select id from %s where %s='%s')", table, field, value)
	fmt.Println("Check Ex : ", query)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var res bool
	for rows.Next() {
		err := rows.Scan(&res)
		if err != nil {
			return false, err
		}
	}
	return res, nil
}
