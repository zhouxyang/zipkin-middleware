package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
	"net/http"
	"time"
)

var tracer *zipkin.Tracer

func IndexHandlerbaidu(w http.ResponseWriter, r *http.Request) {
	client, err := zipkinhttp.NewClient(tracer)
	if err != nil {
		fmt.Fprintln(w, "NewClient err", err)
		return
	}
	req, err := http.NewRequest("GET", "http://localhost:8000/sina", nil)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	/*
		resp, err := client.DoWithAppSpan(req, "requestid")
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	*/
	req = req.WithContext(r.Context())
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer resp.Body.Close()
	fmt.Fprintln(w, "success")
}
func IndexHandlersina(w http.ResponseWriter, r *http.Request) {
	client, err := zipkinhttp.NewClient(tracer)
	if err != nil {
		fmt.Fprintln(w, "NewClient err", err)
		return
	}
	req, err := http.NewRequest("GET", "http://localhost:8000/sohu", nil)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	req = req.WithContext(r.Context())
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	span, _ := tracer.StartSpanFromContext(r.Context(), "some_operation")
	time.Sleep(time.Millisecond)
	span.Finish()
	defer resp.Body.Close()
	fmt.Fprintln(w, "success")
}

func IndexHandlersohu(w http.ResponseWriter, r *http.Request) {
	client, err := zipkinhttp.NewClient(tracer)
	if err != nil {
		fmt.Fprintln(w, "NewClient err", err)
		return
	}
	req, err := http.NewRequest("GET", "http://www.sohu.com/zzz", nil)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	req = req.WithContext(r.Context())
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer resp.Body.Close()
	fmt.Fprintln(w, "success")
}
func main() {
	// create a reporter to be used by the tracer
	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()
	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("myService", "localhost:0")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	tracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)
	r := mux.NewRouter()
	r.Handle("/baidu", alice.New(serverMiddleware).Then(http.HandlerFunc(IndexHandlerbaidu)))
	r.Handle("/sina", alice.New(serverMiddleware).Then(http.HandlerFunc(IndexHandlersina)))
	r.Handle("/sohu", alice.New(serverMiddleware).Then(http.HandlerFunc(IndexHandlersohu)))
	http.ListenAndServe(":8000", r)
}
