package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

var version string

type ClientPrometheus struct {
	Port      string
	PortInner string
	Msgs      []MsgTrack
	metrics   *Metrics

	exchanger chan MsgTrack
}

func NewClientPrometheus(e chan MsgTrack) *ClientPrometheus {
	var p ClientPrometheus
	p.init()
	p.exchanger = e
	p.Msgs = append(p.Msgs, MsgTrack{
		MsgTopic:  "Start INFO",
		MsgTime:   -1,
		MsgOffset: -1,
	})
	return &p
}

func (p *ClientPrometheus) init() {
	version = "1.0.0"
	p.Port = os.Getenv("PROMETHEUS_PORT")
	p.PortInner = os.Getenv("INNER_PORT")
}

func (p *ClientPrometheus) StartHandling() {
	reg := prometheus.NewRegistry()
	p.metrics = NewMetrics(reg)

	p.metrics.Messages.Set(float64(len(p.Msgs)))
	p.metrics.Time.Set(float64(0))
	p.metrics.info.With(prometheus.Labels{"version": version}).Set(1)

	messagesMux := mux.NewRouter()
	messagesMux.HandleFunc("/messages", p.registerMessages)

	metricsMux := mux.NewRouter()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	metricsMux.Handle("/metrics", promHandler)

	log.Println("------Starting Servers------")
	go func() {
		log.Fatal(http.ListenAndServe(p.PortInner, messagesMux))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(p.Port, metricsMux))
	}()
}

func (p *ClientPrometheus) SendMessage(m MsgTrack) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST",
		fmt.Sprintf("http://consumer%v/messages", p.PortInner), bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("[ERR] can't send data :", err)
	}
	defer func() { _ = response.Body.Close() }()
}
