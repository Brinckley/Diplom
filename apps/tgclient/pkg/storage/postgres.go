package storage

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"strings"
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
	dsn := p.cDsnL + p.cDBPassword + p.cDsnR

	conn, err := pgx.Connect(context.Background(), dsn)
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
	dsn := p.cDsnL + p.cDBPassword + p.cDsnR

	conn, err := pgx.Connect(context.Background(), dsn)
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
	dsn := p.cDsnL + p.cDBPassword + p.cDsnR

	conn, err := pgx.Connect(context.Background(), dsn)
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

func (p *TgPostgres) Subscribe(message *tgbotapi.Message) (bool, error) {
	// 1. check artist existence
	// 2. check user existence
	// 3. check pair existence
	// 4. add pair if new
	artistName := strings.TrimSpace(message.Text)
	dsn := p.cDsnL + p.cDBPassword + p.cDsnR

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return false, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	if existence, err := checkExistence(conn, p.cTableUsers, "username", message.From.UserName); err != nil || existence == false {
		if existence == true {
			return false, ErrNoSuchUser
		}
		return false, err
	}
	if existence, err := checkExistence(conn, p.cTableArtists, "name", artistName); err != nil || existence == false {
		if existence == true {
			return false, err
		}
		return false, ErrNoArtists
	}

	subscribe, err := p.checkAndSubscribe(message.From.UserName, artistName)
	if err != nil {
		return false, err
	}
	return subscribe, nil
}

func (p *TgPostgres) Unsubscribe(message *tgbotapi.Message) (bool, error) {
	// 1. check artist existence
	// 2. check user existence
	// 3. check pair existence
	// 4. remove pair if exists
	artistName := strings.TrimSpace(message.Text)
	dsn := p.cDsnL + p.cDBPassword + p.cDsnR

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return false, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	if existence, err := checkExistence(conn, p.cTableUsers, "username", message.From.UserName); err != nil || existence == false {
		if existence == true {
			return false, ErrNoSuchUser
		}
		return false, err
	}
	if existence, err := checkExistence(conn, p.cTableArtists, "name", artistName); err != nil || existence == false {
		if existence == true {
			return false, ErrNoArtists
		}
		return false, err
	}

	subscribe, err := p.checkAndUnsubscribe(message.From.UserName, artistName)
	if err != nil {
		return false, err
	}
	return subscribe, nil
}

func (p *TgPostgres) checkAndSubscribe(username string, artist string) (bool, error) {
	dsn := p.cDsnL + p.cDBPassword + p.cDsnR
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return false, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()
	userId, err := p.getId(conn, p.cTableUsers, "username", username)
	if err != nil {
		return false, fmt.Errorf("error connecting to table %s : %s", p.cTableUsers, err)
	}
	artistId, err := p.getId(conn, p.cTableArtists, "name", artist)
	if err != nil {
		return false, fmt.Errorf("error connecting to table %s : %s", p.cTableArtists, err)
	}
	subscription, err := p.checkSubscription(conn, userId, artistId)
	if err != nil {
		return false, err
	}
	if subscription {
		return false, nil
	}

	queryInsert := fmt.Sprintf("insert into %s (user_id, artist_id) "+
		"values ($1, $2) returning id;", p.cTableUserArtist)
	fmt.Println("Insert q : ", queryInsert)
	var pairId int
	err = conn.QueryRow(context.Background(),
		queryInsert,
		userId, artistId).Scan(&pairId)
	if err != nil {
		return false, err
	}
	log.Println("Pair added with id :", pairId)

	return true, nil
}

func (p *TgPostgres) checkAndUnsubscribe(username string, artist string) (bool, error) {
	dsn := p.cDsnL + p.cDBPassword + p.cDsnR
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return false, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()
	userId, err := p.getId(conn, p.cTableUsers, "username", username)
	if err != nil {
		return false, fmt.Errorf("error connecting to table %s : %s", p.cTableUsers, err)
	}
	artistId, err := p.getId(conn, p.cTableArtists, "name", artist)
	if err != nil {
		return false, fmt.Errorf("error connecting to table %s : %s", p.cTableArtists, err)
	}
	subscription, err := p.checkSubscription(conn, userId, artistId)
	if err != nil {
		return false, err
	}
	if !subscription {
		return false, nil
	}

	queryDelete := fmt.Sprintf("delete from %s where user_id=%v and artist_id=%v returning id",
		p.cTableUserArtist, userId, artistId)
	fmt.Println("Delete q : ", queryDelete)
	var pairId int
	err = conn.QueryRow(context.Background(),
		queryDelete).Scan(&pairId)
	if err != nil {
		return false, err
	}
	log.Println("Pair deleted with id :", pairId)

	return true, nil
}

func (p *TgPostgres) checkSubscription(conn *pgx.Conn, uId int, aId int) (bool, error) {
	queryGetId := fmt.Sprintf("select id from %s where user_id=%v and artist_id=%v",
		p.cTableUserArtist, uId, aId)
	fmt.Println("Check Subs : ", queryGetId)
	rows, err := conn.Query(context.Background(), queryGetId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	res := -1
	for rows.Next() {
		err := rows.Scan(&res)
		if err != nil {
			return false, err
		}
	}
	if res == -1 {
		return false, nil
	}
	return true, nil
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
