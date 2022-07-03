package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("Received request for hello world")
	io.WriteString(w, "Hello World\n")
}

func metricsMiddlewareBuilder(reqCounter prometheus.Counter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			reqCounter.Inc()
		})
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cur := time.Now()
		log.Printf("Starting request for %s", req.URL.Path)
		next.ServeHTTP(w, req)
		diff := time.Since(cur)
		log.Printf("Completed request for %s in %d milliseconds", req.URL.Path, diff.Milliseconds())
	})
}

func applyMiddlewares(mux *http.ServeMux) {
	requestsServed := promauto.NewCounter(prometheus.CounterOpts{
		Name: "promexperiment_processed_requests_total",
		Help: "The total number of processed requests",
	})
	metricsMiddleware := metricsMiddlewareBuilder(requestsServed)

	mux.Handle("/hello", loggingMiddleware(metricsMiddleware(http.HandlerFunc(helloHandler))))
	mux.Handle("/metrics", promhttp.Handler())
}

func main() {
	mux := http.NewServeMux()

	applyMiddlewares(mux)

	log.Print("Listening on :8080...")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
