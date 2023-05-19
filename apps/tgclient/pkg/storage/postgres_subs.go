package storage

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v5"
	"log"
	"strings"
	"time"
)

func (p *TgPostgres) Subscribe(message *tgbotapi.Message) (bool, error) {
	// 1. check artist existence
	// 2. check user existence
	// 3. check pair existence
	// 4. add pair if new
	artistName := strings.TrimSpace(message.Text)
	conn, err := pgx.Connect(context.Background(), p.cDsn)
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
	conn, err := pgx.Connect(context.Background(), p.cDsn)
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

func (p *TgPostgres) updateSubscriptionTimeStamp(subscribers []int64) error {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	for _, s := range subscribers {
		updTimeStampQuery := fmt.Sprintf(
			"update %s set updTime = %v where %s.user_id in "+
				"(select id from %s where %s.сhatID = %v);",
			p.cTableUserArtist, s, p.cTableUserArtist,
			p.cTableUsers, p.cTableUsers, s)
		fmt.Println("[QUERY @] ", updTimeStampQuery)
		_, err = conn.Exec(context.Background(), updTimeStampQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *TgPostgres) checkAndSubscribe(username string, artist string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
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

	queryInsert := fmt.Sprintf("insert into %s (user_id, artist_id, updTime) "+
		"values ($1, $2, $3) returning id;", p.cTableUserArtist)
	fmt.Println("Insert q : ", queryInsert)
	var pairId int
	err = conn.QueryRow(context.Background(),
		queryInsert,
		userId, artistId, time.Now().Unix()).Scan(&pairId)
	if err != nil {
		return false, err
	}
	log.Println("Pair added with id :", pairId)

	return true, nil
}

func (p *TgPostgres) checkAndUnsubscribe(username string, artist string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), p.cDsn)
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
