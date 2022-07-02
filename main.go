package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "promexperiment_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("Received request for hello world")
	io.WriteString(w, "Hello World\n")
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
	mux.Handle("/hello", loggingMiddleware(http.HandlerFunc(helloHandler)))
}

func main() {
	mux := http.NewServeMux()

	applyMiddlewares(mux)

	log.Print("Listening on :8080...")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
