package utils

import (
	"strings"
)

const (
	classic    = "классика"           // классика
	rock       = "рок"                // рок
	pop        = "поп"                // поп, эстрада
	jazz       = "джаз"               // джаз
	folk       = "фолк"               // народная, фолк
	chanson    = "шансон"             // авторская, шансон, романсы
	rap        = "рэп"                // хип-хоп, рэп
	electronic = "электронная музыка" // электронная музыка
)

var classicDesc = []string{"классика", "classical", "neoclassical"}
var rockDesc = []string{"рок", "rock", "metal"}
var popDesc = []string{"поп", "эстрада", "pop", "russian pop", "retro"}
var jazzDesc = []string{"джаз", "jazz"}
var folkDesc = []string{"народная", "фолк", "folk"}
var chansonDesc = []string{"авторская", "шансон", "романсы", "bard", "chanson", "romance"}
var rapDesc = []string{"хип-хоп", "рэп", "hip-hop", "rap"}
var electronicDesc = []string{"электронная музыка", "electronic", "electro"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetTag(tag string) string {
	genre := strings.ToLower(tag)
	if contains(classicDesc, genre) {
		return classic
	}
	if contains(rockDesc, genre) {
		return rock
	}
	if contains(popDesc, genre) {
		return pop
	}
	if contains(jazzDesc, genre) {
		return jazz
	}
	if contains(folkDesc, genre) {
		return folk
	}
	if contains(chansonDesc, genre) {
		return chanson
	}
	if contains(rapDesc, genre) {
		return rap
	}
	if contains(electronicDesc, genre) {
		return electronic
	}
	return ""
}
