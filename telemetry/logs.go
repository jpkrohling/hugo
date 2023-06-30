package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var logger *zap.Logger

func NewLogger() (*zap.Logger, error) {
	var err error
	logger, err = zap.NewDevelopment()
	return logger, err
}

func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	sp := trace.SpanFromContext(ctx)

	traceID := sp.SpanContext().TraceID().String()
	spanID := sp.SpanContext().SpanID().String()

	fields = append(fields, zap.String("traceID", traceID), zap.String("spanID", spanID))

	logger.Info(msg, fields...)
}
