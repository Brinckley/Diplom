package consul

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
	"net/http"
	kafka_producer "producer/internal/kafka-producer"
)

func consulCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "consulCheck")
}

func RegisterServer() {
	config := consul.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consul.NewClient(config)

	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	checkPort := kafka_producer.KafkaPort

	registration := new(consul.AgentServiceRegistration)
	registration.ID = "producerNode_1"
	registration.Name = "producerNode"
	registration.Port = 9527
	registration.Tags = []string{"producerNode"}
	registration.Address = "127.0.0.1"
	registration.Check = &consul.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s", // CHECK deletes this service 30 seconds after failure
	}

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Fatal("register server error : ", err)
	}

	http.HandleFunc("/check", consulCheck)
	http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)

}
