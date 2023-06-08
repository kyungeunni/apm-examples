package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func write(writer http.ResponseWriter, message string) {
	_, err := writer.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

// var tracingClient = apmhttp.WrapClient(http.DefaultClient)

func englishHandler(writer http.ResponseWriter, request *http.Request) {
	// TODO: make outgoing http requests to add more spans
	resp, err := otelhttp.Get(request.Context(), "https://kyungeun.kim")
	if err != nil {
		// apm.CaptureError(request.Context(), err).Send()
		// span.SetStatus(codes.Error, "operationThatCouldFail failed")
		// span.RecordError(err)
		http.Error(writer, "failed to query", 500)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// apm.CaptureError(request.Context(), err).Send()
		http.Error(writer, "failed to parse body", 500)
		return
	}
	write(writer, string(body[:]))
}

func frenchHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	_, span := trace.SpanFromContext(request.Context()).TracerProvider().Tracer("exampleTracer").Start(ctx, "Delay") // , "Delay", "request")
	defer span.End()
	time.Sleep(3 * time.Second)
	span.SetAttributes(attribute.String("reason", "time.Sleep"))
	write(writer, "Salut, web!")
}

func hindiHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Namaste, web!")
}

func main() {
	ctx := context.Background()
	tp, err := newTraceProvider(ctx)

	if err != nil {
		log.Fatal(err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()
	mux := http.NewServeMux()
	mux.Handle("/hello", otelhttp.NewHandler(http.HandlerFunc(englishHandler), "hello"))
	mux.Handle("/salut", otelhttp.NewHandler(http.HandlerFunc(frenchHandler), "salut"))
	mux.Handle("/namaste", otelhttp.NewHandler(http.HandlerFunc(hindiHandler), "namaste"))
	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return &otlptrace.Exporter{}, err
	}
	return exporter, nil
}

func newTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exp, err := newExporter(ctx)
	if err != nil {
		return sdktrace.NewTracerProvider(), err
	}
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("my-hello-web-server"),
		),
	)

	if err != nil {
		return sdktrace.NewTracerProvider(), nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
