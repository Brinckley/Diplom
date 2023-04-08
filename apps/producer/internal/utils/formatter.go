package utils

import (
	"log"
	"net/url"
)

func DecodeUrl(oldUrl string) string {
	decodedValue, err := url.QueryUnescape(oldUrl)
	if err != nil {
		log.Printf("[ERR] can't decode the url : %s, The URL is : '%s'", err.Error(), oldUrl)
		return ""
	}
	return decodedValue
}
