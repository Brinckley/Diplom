package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	_, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": []string{"localhost"},
		"group.id":          "artist-gr",
		"auto.offset.reset": "earliest",
	})
	fmt.Println(err)
}
