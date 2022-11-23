package services

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/grpclog"
	"log"
	"os"
	"time"
)

var cEndPoint = ""

func InitETCD() {
	cEndPoint = os.Getenv("END_POINT")
}

func ConnectETCD() {
	clientv3.SetLogger(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"etcd:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln("Error creating new client :", err)
	}
	defer cli.Close() // make sure to close the client

	_, err = cli.Put(context.Background(), "foo", "bar")
	if err != nil {
		log.Fatalln("Error putting kv", err)
	}

	foo, err := cli.Get(context.Background(), "foo")
	if err != nil {
		log.Fatalln("Error getting kv", err)
	}

	log.Println("FOO KVS: ", foo.Kvs)
}
