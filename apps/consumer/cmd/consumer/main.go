package main

import (
	"consumer/internal/kafka"
	"consumer/internal/postgres"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	// init environment
	logger := logrus.Logger{Formatter: &logrus.JSONFormatter{}} // creating and initializing the logger
	logger.SetOutput(os.Stdout)

	postgresClient := postgres.NewPostgres(&logger)
	kafkaClient := kafka.NewKafka(postgresClient, &logger)
	kafkaClient.ConsumeAndSend()
}
