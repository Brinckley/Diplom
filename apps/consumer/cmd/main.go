package main

import (
	"consumer/pkg/kafka"
	"consumer/pkg/kafka/prometheus"
	"consumer/pkg/postgres"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	// init all environment

	prom := prometheus.NewClientPrometheus()

	logger := logrus.Logger{Formatter: &logrus.JSONFormatter{}} // creating and initializing the logger
	logger.SetOutput(os.Stdout)

	postgresClient := postgres.NewPostgres(&logger)
	clientKafka := kafka.NewKafka(postgresClient, &logger, prom)
	//clientKafka := kafka.NewKafka(nil, &logger, prom, ) // postgres taken away for debug

	clientKafka.ConsumeAndSend()
}
