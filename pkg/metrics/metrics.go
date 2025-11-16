package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    EmailsSent = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "email_sent_total",
        Help: "Total number of successfully sent emails",
    })

	EmailsFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "email_failed_total",
		Help: "Total number of failed emails",
	})

    EmailDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:    "email_send_duration_seconds",
        Help:    "Time taken to send an email",
        Buckets: prometheus.DefBuckets,
    })

	WorkerActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "email_worker_active_count",
        Help: "Number of active workers",
	})

)

func Init(){
	prometheus.MustRegister(EmailsSent)
	prometheus.MustRegister(EmailsFailed)
	prometheus.MustRegister(EmailDuration)
	prometheus.MustRegister(WorkerActive)
}

func StartMetrics(addr string) {

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("[prometheus] starting metrics server at %s/metrics", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("metrics server failed: %v", err)
	}
	
}