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

func main() {
	kafka_producer.InitProducer()
	consulDebug := false

	if !consulDebug {
		Parsing()
	} else {
		consul.RegisterServer()
	}

}

func Parsing() {
	inputPath := os.Getenv("INPUT_PATH")
	//outputPath := os.Getenv("OUTPUT_PATH")
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
