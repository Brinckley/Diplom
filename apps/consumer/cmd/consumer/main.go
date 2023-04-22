package main

import (
	"consumer/internal/kafka"
	"github.com/sirupsen/logrus"
	"os"
)

// https://docs.kudago.com/api/#

func main() {
	// init all environment
	logger := logrus.Logger{Formatter: &logrus.JSONFormatter{}} // creating and initializing the logger
	logger.SetOutput(os.Stdout)

	//postgresClient := postgres.NewPostgres(&logger)
	//clientKafka := kafka.NewKafka(postgresClient, &logger)
	clientKafka := kafka.NewKafka(nil, &logger)

	clientKafka.ConsumeAndSend()
}
