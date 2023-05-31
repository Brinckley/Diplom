package prometheus

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
)

type MsgTrack struct {
	MsgTopic  string `json:"msgTopic"`
	MsgTime   int64  `json:"msgTime"`
	MsgOffset int64  `json:"msgOffset"`
}

func NewMsgTrack(message kafka.Message, delta int64) MsgTrack {
	msg := MsgTrack{
		MsgTopic:  message.Topic,
		MsgTime:   delta,
		MsgOffset: message.Offset,
	}
	return msg
}

type Metrics struct {
	Messages prometheus.Gauge
	info     *prometheus.GaugeVec
	Time     prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "Consumer",
			Name:      "consumer_version",
			Help:      "Information about the Consumer version.",
		},
			[]string{"version"}),
		Messages: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Consumer",
			Name:      "read_messages",
			Help:      "Number of read messages.",
		}),
		Time: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Consumer",
			Name:      "time_read",
			Help:      "Time for reading message.",
		}),
	}
	reg.MustRegister(m.info, m.Messages, m.Time)
	return m
}

func (p *ClientPrometheus) registerMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p.getMsgMetrics(w, r)
	case "POST":
		p.createMsgMetrics(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (p *ClientPrometheus) getMsgMetrics(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(p.Msgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (p *ClientPrometheus) createMsgMetrics(w http.ResponseWriter, r *http.Request) {
	var mt MsgTrack
	err := json.NewDecoder(r.Body).Decode(&mt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.Msgs = append(p.Msgs, mt)
	log.Println("Num : ", len(p.Msgs))
	log.Println("Time : ", mt.MsgTime)
	log.Println("Topic : ", mt.MsgTopic)
	log.Println("Offset : ", mt.MsgOffset)
	log.Println("**************")
	p.metrics.Messages.Set(float64(len(p.Msgs)))
	p.metrics.Time.Set(float64(mt.MsgTime))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message added!"))
}
