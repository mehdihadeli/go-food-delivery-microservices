package metrics

import (
	"fmt"
	metric2 "go.opentelemetry.io/otel/metric"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func AddOtelMetrics() metric2.Meter {
	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter := otelprom.New()
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("github.com/open-telemetry/opentelemetry-go/example/prometheus")

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics(exporter.Collector)

	return meter
}

func serveMetrics(collector prometheus.Collector) {
	registry := prometheus.NewRegistry()
	err := registry.Register(collector)
	if err != nil {
		fmt.Printf("error registering collector: %v", err)
		return
	}

	log.Printf("serving metrics at localhost:2222/metrics")
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	err = http.ListenAndServe(":2222", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
