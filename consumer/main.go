package main

import (
	"consumer/services"
)

func main() {
	services.InitConsumer()
	services.InitDatabase()
	services.InitETCD()
	//dsnL := os.Getenv("DSN_LEFT")
	//dsnR := os.Getenv("DSN_RIGHT")
	//services.PGXInsert(dsnL, dsnR)
	services.ConnectETCD()
}
