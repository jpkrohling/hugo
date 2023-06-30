package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/samplers/jaegerremote"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var cl = otlptracegrpc.NewClient()

func NewTracerProvider() (*sdktrace.TracerProvider, error) {
	exp, err := otlptrace.New(context.Background(), cl)
	if err != nil {
		return nil, err
	}

	_, err = stdouttrace.New()
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			b3.New(),
		),
	)

	_ = jaegerremote.New(
		"hugo",
		jaegerremote.WithSamplingServerURL("http://localhost:5778/sampling"),
		jaegerremote.WithSamplingRefreshInterval(10*time.Second),
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(NewResource()),
	), nil

}

func Shutdown(ctx context.Context) error {
	return cl.Stop(ctx)
}

func GetTracer() trace.Tracer {
	return otel.Tracer("hugo")
}
