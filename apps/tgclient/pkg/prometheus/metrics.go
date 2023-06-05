package prometheus

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type EventNotification struct {
	EventTitle string
	EventTime  int64
}

func NewEventNotification(title string, delta int64) EventNotification {
	en := EventNotification{
		EventTitle: title,
		EventTime:  delta,
	}
	return en
}

type Metrics struct {
	info *prometheus.GaugeVec
	Time prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "TgClient",
			Name:      "tgclient_version",
			Help:      "Information about the Telegram Client version.",
		},
			[]string{"version"}),
		Time: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "TgClient",
			Name:      "time_notify",
			Help:      "Time for reading and sending notifications.",
		}),
	}
	reg.MustRegister(m.info, m.Time)
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
	b, err := json.Marshal(p.EvNotification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	func() { _, _ = w.Write(b) }()
}

func (p *ClientPrometheus) createMsgMetrics(w http.ResponseWriter, r *http.Request) {
	var en EventNotification
	err := json.NewDecoder(r.Body).Decode(&en)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.metrics.Time.Set(float64(en.EventTime))
	w.WriteHeader(http.StatusCreated)
	func() { _, _ = w.Write([]byte("Message added!")) }()
}
