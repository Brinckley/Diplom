package kafka

//
//import (
//	prometheus2 "consumer/pkg/kafka/prometheus"
//	"encoding/json"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"log"
//	"net/http"
//)
//
//func (k *ClientKafka) StartHandling() {
//	reg := prometheus.NewRegistry()
//	m := prometheus2.NewMetrics(reg)
//	m.ArtistTime.Set()
//
//	http.Handle("/metrics", promhttp.Handler())
//	http.HandleFunc("/kafka", k.GetKafkaMetricsHandler)
//
//	err := http.ListenAndServe(k.prometheusClient.Port, nil)
//	if err != nil {
//		log.Println("[ERR] can't listen and serve at port ", k.prometheusClient.Port)
//		return
//	}
//}
//
//func (k *ClientKafka) GetKafkaMetricsHandler(w http.ResponseWriter, r *http.Request) {
//	artistStats, err := json.Marshal(k.artistReader.Stats().ReadTime.Avg)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	_, err = w.Write(artistStats)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//}
