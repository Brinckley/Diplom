package kassir_functions

import (
	"fmt"
	"log"
)

func CityAbbr(cityName string) (string, error) {
	fullUrl := fmt.Sprintf("https://msk.kassir.ru/")

	doc, err := getHTMLFromLink(fullUrl)
	if err != nil {
		log.Println("no such url")
		return "", err
	}

	docc := doc.
		Find("div", "class", "modal fade city-select in")
	if docc.Error != nil {
		log.Println("no city info found")
		return "", err
	}

	return "", nil
}

/*
"https://msk.kassir.ru/" Москва
"https://spb.kassir.ru/" Санкт-Петербург
"https://aba.kassir.ru/" Абакан
"https://anapa.kassir.ru/" Анапа
"https://arh.kassir.ru/" Архангельск
"https://astr.kassir.ru/" Астрахань
"https://brn.kassir.ru/" Барнаул
"https://belgorod.kassir.ru/" Белгород
"https://blag.kassir.ru/" Благовещенск
"https://bryansk.kassir.ru/" Брянск
"https://nov.kassir.ru/"  Великий Новгород
"https://vl.kassir.ru/"  Владивосток
"https://vlm.kassir.ru/"  Владимир
"https://vlg.kassir.ru/"  Волгоград
"https://vologda.kassir.ru/"  Вологда
"https://vrn.kassir.ru/"  Воронеж
"https://gel.kassir.ru/"  Геленджик
"https://ekb.kassir.ru/"  Екатеринбург
"https://ivanovo.kassir.ru/"  Иваново
"https://izhevsk.kassir.ru/"  Ижевск
"https://irk.kassir.ru/"  Иркутск
"https://yola.kassir.ru/"  Йошкар-Ола
"https://kzn.kassir.ru/"  Казань
"https://kgd.kassir.ru/"  Калининград
"https://klg.kassir.ru/"  Калуга
"https://kemerovo.kassir.ru/"  Кемерово
"https://kirov.kassir.ru/"  Киров
"https://komsomolsk.kassir.ru/"  Комсомольск-на-Амуре
"https://krd.kassir.ru/"  Краснодар
"https://krs.kassir.ru/"  Красноярск
"https://kursk.kassir.ru/"  Курск
"https://lzr.kassir.ru/"  Лазаревское
"https://lipetsk.kassir.ru/"  Липецк
"https://mgn.kassir.ru/"  Магнитогорск
"https://murm.kassir.ru/"  Мурманск
"https://nabchelny.kassir.ru/"  Набережные Челны
"https://nn.kassir.ru/"  Нижний Новгород
"https://novokuznetsk.kassir.ru/"  Новокузнецк
"https://nvrsk.kassir.ru/"  Новороссийск
"https://nsk.kassir.ru/"  Новосибирск
"https://omsk.kassir.ru/"  Омск
"https://orel.kassir.ru/"  Орёл
"https://orenburg.kassir.ru/"  Оренбург
"https://orsk.kassir.ru/"  Орск
"https://pnz.kassir.ru/"  Пенза
"https://perm.kassir.ru/"  Пермь
"https://ptz.kassir.ru/"  Петрозаводск
"https://pskov.kassir.ru/"  Псков
"https://rnd.kassir.ru/"  Ростов-на-Дону
"https://rzn.kassir.ru/"  Рязань
"https://smr.kassir.ru/"  Самара
"https://saransk.kassir.ru/"  Саранск
"https://saratov.kassir.ru/"  Саратов
"https://smolensk.kassir.ru/"  Смоленск
"https://sochi.kassir.ru/"  Сочи
"https://sk.kassir.ru/"  Ставрополь
"https://oskol.kassir.ru/"  Старый Оскол
"https://sur.kassir.ru/"  Сургут
"https://tambov.kassir.ru/"  Тамбов
"https://tver.kassir.ru/"  Тверь
"https://tlt.kassir.ru/"  Тольятти
"https://tomsk.kassir.ru/"  Томск
"https://tula.kassir.ru/"  Тула
"https://tmn.kassir.ru/"  Тюмень
"https://ulan.kassir.ru/"  Улан-Удэ
"https://ulyanovsk.kassir.ru/"  Ульяновск
"https://ufa.kassir.ru/"  Уфа
"https://hbr.kassir.ru/"  Хабаровск
"https://chaik.kassir.ru/"  Чайковский
"https://cheboksary.kassir.ru/"  Чебоксары
"https://chel.kassir.ru/"  Челябинск
"https://cher.kassir.ru/"  Череповец
"https://chita.kassir.ru/"  Чита
"https://sakh.kassir.ru/"  Южно-Сахалинск
"https://yar.kassir.ru/"  Ярославль
*/
