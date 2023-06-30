package telemetry

import (
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var FooCounter metric.Int64Counter

func NewMeterProvider() (metric.MeterProvider, error) {
	exp, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp, sdkmetric.WithInterval(time.Second))),
		sdkmetric.WithResource(NewResource()),
	)

	m := provider.Meter("hugo")

	FooCounter, err = m.Int64Counter(
		"io.gohugoio.server.foo",
		metric.WithDescription("Mede o numero de foos."),
	)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func GetMeter() metric.Meter {
	return otel.Meter("hugo")
}
