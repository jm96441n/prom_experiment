package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"net/http"
	"promexperiment/webapi"
)

func main() {
	logger := &log.Logger{}
	requestsServed := promauto.NewCounter(prometheus.CounterOpts{
		Name: "promexperiment_processed_requests_total",
		Help: "The total number of processed requests",
	})

	mux := http.NewServeMux()

	webapi.SetupServer(logger, mux, requestsServed)

	log.Print("Listening on :8080...")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
