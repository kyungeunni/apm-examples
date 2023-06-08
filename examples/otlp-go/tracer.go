// package main

// import (
// 	"context"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
// 	"go.opentelemetry.io/otel/trace"
// )

// var tracer trace.Tracer

// func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
// 	client := otlptracehttp.NewClient()
// 	exporter, err := otlptrace.New(ctx, client)
// 	if err != nil {
// 		return &otlptrace.Exporter{}, err
// 	}
// 	return exporter, nil
// }

// func newTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
// 	exp, err := newExporter(ctx)
// 	if err != nil {
// 		return sdktrace.NewTracerProvider(), err
// 	}
// 	// Ensure default SDK resources and the required service name are set.
// 	r, err := resource.Merge(
// 		resource.Default(),
// 		resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			semconv.ServiceName("my-hello-web-server"),
// 		),
// 	)

// 	if err != nil {
// 		return sdktrace.NewTracerProvider(), nil
// 	}

// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithBatcher(exp),
// 		sdktrace.WithResource(r),
// 	)

// 	otel.SetTracerProvider(tp)

// 	return tp, nil
// }
