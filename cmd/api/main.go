package main

import (
	"context"
	"llm-agent-go/cmd/api/routes"
	"llm-agent-go/cmd/service_container"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

func main() {
	ctx := context.Background()

	container := service_container.NewServiceContainer()
	traceProvider, err := initTraces(ctx, "llm-agent-go")
	if err != nil {
		container.Logger.Fatal().Err(err).Msg("failed to init traces")
	}
	defer traceProvider.Shutdown(ctx)

	routes.InitRoutes(container.Controllers)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initTraces(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("alloy:4317"),
		otlptracegrpc.WithTimeout(5*time.Second),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	r, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}
