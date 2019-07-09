// +build go1.7

package main

import (
	"fmt"
	"net/http"
	"os"

	opentracing "github.com/opentracing/opentracing-go"

	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin-contrib/zipkin-go-opentracing/examples/cli_with_2_services/svc1"
	"github.com/openzipkin-contrib/zipkin-go-opentracing/examples/cli_with_2_services/svc2"
)

const (
	// Our service name.
	serviceName = "svc1"

	// Host + port of our service.
	hostPort = "127.0.0.1:61001"

	// Endpoint to send Zipkin spans to.
	zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"

	// Debug mode.
	debug = false

	// Base endpoint of our SVC2 service.
	svc2Endpoint = "http://localhost:61002"

	// same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
	sameSpan = true

	// make Tracer generate 128 bit traceID's for root spans.
	traceID128Bit = true
)

//svc1
func main() {
	// create collector.
	collector, err := zipkin.NewHTTPCollector(zipkinHTTPEndpoint)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// create recorder.
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)

	// create tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// explicitly set our tracer to be the default tracer.
	opentracing.InitGlobalTracer(tracer)

	// create the client to svc2
	svc2Client := svc2.NewHTTPClient(tracer, svc2Endpoint)

	// create the service implementation
	service := svc1.NewService(svc2Client)

	// create the HTTP Server Handler for the service
	handler := svc1.NewHTTPHandler(tracer, service)

	// start the service
	fmt.Printf("Starting %s on %s\n", serviceName, hostPort)
	http.ListenAndServe(hostPort, handler)
}
