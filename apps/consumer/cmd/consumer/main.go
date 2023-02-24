package main

import (
	"consumer/internal/etcd"
	"consumer/internal/kafka"
	"consumer/internal/postgres"
)

// https://docs.kudago.com/api/#

func main() {
	// init all environment
	kafka.InitConsumer()
	postgres.InitDatabase()
	etcd.InitETCD()

	kafka.ConsumeAndSend()

	// consuming data from three kafka topics
	//dAr, dAl, dTr := kafka.ReadAll(context.Background(), 0)

	// inserting data into db
	//postgres.PGXInsert(dAr, dAl, dTr)

	// connecting to etcd
	// etcd.ConnectETCD()
}
