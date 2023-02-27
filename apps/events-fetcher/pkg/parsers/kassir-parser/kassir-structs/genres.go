package kassir_structs

import (
	"errors"
	"strings"
)

const (
	classic    = "category%5B%5D=3001" // классика
	rock       = "category%5B%5D=3002" // рок
	pop        = "category%5B%5D=3003" // поп, эстрада
	jazz       = "category%5B%5D=3004" // джаз
	folk       = "category%5B%5D=3005" // народная, фолк
	chanson    = "category%5B%5D=3006" // авторская, шансон, романсы
	rap        = "category%5B%5D=3007" // хип-хоп, рэп
	electronic = "category%5B%5D=3008" // электронная музыка
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

func SelectGenre(genre string) (string, error) {
	genre = strings.ToLower(genre)
	if contains(classicDesc, genre) {
		return classic, nil
	}
	if contains(rockDesc, genre) {
		return rock + "&" + classic, nil
	}
	if contains(popDesc, genre) {
		return pop, nil
	}
	if contains(jazzDesc, genre) {
		return jazz, nil
	}
	if contains(folkDesc, genre) {
		return folk, nil
	}
	if contains(chansonDesc, genre) {
		return chanson, nil
	}
	if contains(rapDesc, genre) {
		return rap, nil
	}
	if contains(electronicDesc, genre) {
		return electronic, nil
	}
	return "", errors.New("error recognizing genre")
}
