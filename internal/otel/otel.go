package otel

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/richardbizik/gommentary/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func SetupOtel(conf config.Config) (http.Handler, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	resource, err := resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceName("gommentary"),
		attribute.String("library.language", "go"),
	))
	if err != nil {
		return nil, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(resource),
	)
	otel.SetMeterProvider(meterProvider)
	return promhttp.Handler(), nil
}
