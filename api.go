package main

import (
	"demo/my"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporter/trace/stackdriver"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var projectID = "gitlab-203909"

func main() {
	initTracer()

	r := mux.NewRouter()
	r.HandleFunc("/", my.MainHandler)

	http.ListenAndServe(":8080", r)
}

func initTracer() {
	// Create Stackdriver exporter to be able to retrieve
	// the collected spans.
	exporter, err := stackdriver.NewExporter(
		stackdriver.WithProjectID(projectID),
	)
	if err != nil {
		log.Fatal(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}
