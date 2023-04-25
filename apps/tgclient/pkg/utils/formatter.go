package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/url"
	"strconv"
	"strings"
)

func DecodeUrl(oldUrl string) string {
	decodedValue, err := url.QueryUnescape(oldUrl)
	if err != nil {
		log.Printf("[ERR] can't decode the url : %s, The URL is : '%s'", err.Error(), oldUrl)
		return ""
	}
	strings.Replace(decodedValue, " ", "%20", -1)
	return decodedValue
}

func ListOfAlbumsFormatter(albums []string) string {
	output := ""
	for id, a := range albums {
		output += strconv.Itoa(id+1) + ". " + a + "\n"
	}

	return output
}

func CallbackToMsg(query *tgbotapi.CallbackQuery) tgbotapi.Message {
	a := []tgbotapi.MessageEntity{{
		Type:   "bot_command",
		Offset: 0,
		Length: len(query.Data),
	}}

	m := tgbotapi.Message{
		MessageID: 0,
		From:      query.From,
		Date:      0,
		Entities:  &a,
		Chat:      query.Message.Chat,
		Text:      query.Data,
	}
	return m
}
