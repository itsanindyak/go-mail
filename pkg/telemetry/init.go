package telemetry

import (
	"context"
	"log"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	grpcLog "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	grpcMetrics "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	grpcTrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitOTel(ctx context.Context) func(context.Context) error {
	exporterUrl := os.Getenv("EXPORTER_URL")

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// semconv.ServiceName("email-engine"),
			semconv.ServiceNameKey.String("email-engine"),
		))

	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// 1. Traces (Tempo)
	traceExp, _ := grpcTrace.New(
		ctx,
		grpcTrace.WithInsecure(),
		grpcTrace.WithEndpoint(exporterUrl),
	)

	traceP := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(traceP)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// 2. metrics
	metricsExp, _ := grpcMetrics.New(
		ctx,
		grpcMetrics.WithInsecure(),
		grpcMetrics.WithEndpoint(exporterUrl),
	)

	metricP := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricsExp)),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(metricP)

	// 3. logs
	logExp, _ := grpcLog.New(
		ctx,
		grpcLog.WithInsecure(),
		grpcLog.WithEndpoint(exporterUrl),
	)

	logP := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
		sdklog.WithResource(res),
	)
	global.SetLoggerProvider(logP)

	logger := otelslog.NewLogger("email-engine")
	slog.SetDefault(logger)

	return func(ctx context.Context) error {
		traceP.Shutdown(ctx)
		metricP.Shutdown(ctx)
		logP.Shutdown(ctx)

		return nil
	}

}
