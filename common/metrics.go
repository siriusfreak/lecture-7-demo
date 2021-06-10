package common

import "github.com/prometheus/client_golang/prometheus"

var processedByHandler = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "processed_by_handler", // metric name
		Help: "Number of processed by handler",
	},
	[]string{"handler"}, // labels
)
var consumedMessages = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "consumed_messages", // metric name
		Help: "Number of consumed messages",
	},
)

func RegisterMetrics() {
	prometheus.MustRegister(processedByHandler)
	prometheus.MustRegister(consumedMessages)
}

func IncProcessedByHandler(handler string) {
	processedByHandler.With(prometheus.Labels{"handler": handler}).Inc()
}

func IncConsumedMessages() {
	consumedMessages.Inc()
}