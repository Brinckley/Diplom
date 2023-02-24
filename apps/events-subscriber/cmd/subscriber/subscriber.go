package main

import (
	"events-subscriber/internal/redis"
	"log"
)

func main() {
	redis.InitRedis()

	client, err := redis.GetSubscriber()
	if err != nil {
		log.Fatalln("error while subscribing occurred : ", err)
	}

	artists := []string{"ддт", "филипп%20киркоров", "константин%20никольский", "nazareth", "metallica", "мельница", "алиса"}

	for _, a := range artists {
		err = redis.GetSubscriptionData(client, a)
		if err != nil {
			log.Printf("error occurred while getting data from %s : %s\n", a, err.Error())
			continue
		}
	}
}
