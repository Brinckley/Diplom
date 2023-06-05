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
	"time"
)

var version string

type ClientPrometheus struct {
	Port           string
	PortInner      string
	EvNotification EventNotification
	metrics        *Metrics
}

func NewClientPrometheus() *ClientPrometheus {
	version = "v.1.0.0"
	var p ClientPrometheus
	p.init()

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

	p.metrics.Time.Set(float64(0))
	p.metrics.info.With(prometheus.Labels{"version": version}).Set(1)

	notificationMux := mux.NewRouter()
	notificationMux.HandleFunc("/notifications", p.registerMessages)

	metricsMux := mux.NewRouter()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	metricsMux.Handle("/metrics", promHandler)

	log.Println("------Starting Servers------")
	go func() {
		log.Fatal(http.ListenAndServe(p.PortInner, notificationMux))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(p.Port, metricsMux))
	}()
}

func (p *ClientPrometheus) SendMessage(m EventNotification) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST",
		fmt.Sprintf("http://tgclient%v/notifications", p.PortInner), bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("[ERR] can't send data :", err)
	}
	defer func() { _ = response.Body.Close() }()
	time.Sleep(1 * time.Second)
}
