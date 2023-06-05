package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"producer/internal/collector"
	"producer/internal/consul"
	"producer/internal/kafka-producer"
)

var DATA_LOAD_DEBUG = false

func main() {
	kafka_producer.InitProducer()
	consulDebug := false

	if !consulDebug {
		if DATA_LOAD_DEBUG {
			CheckLoad(100000)
		} else {
			Parsing()
		}
	} else {
		consul.RegisterServer()
	}

}

func CheckLoad(num int) {
	collector.FalseCollector(num)
}

func Parsing() {
	inputPath := os.Getenv("INPUT_PATH")
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalln("[ERR] can't open the file to read artists : ", err.Error())
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		if name != "" {
			collector.ParserCollectorArtistWithReleases(name)
			fmt.Println(name)
		}
	}
}
