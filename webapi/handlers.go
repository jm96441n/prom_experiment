package webapi

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"time"
)

type incrementer interface {
	Inc()
}

type muxHandler interface {
	Handle(string, http.Handler)
}

type logPrinter interface {
	Print(...any)
	Printf(string, ...any)
}

func helloHandlerBuilder(logger logPrinter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		logger.Print("Received request for hello world")
		io.WriteString(w, "Hello World\n")
	})
}

func counterMiddlewareBuilder(reqCounter incrementer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			reqCounter.Inc()
		})
	}
}

func loggingMiddlewareBuilder(logger logPrinter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			cur := time.Now()
			logger.Printf("Starting request for %s", req.URL.Path)
			next.ServeHTTP(w, req)
			diff := time.Since(cur)
			logger.Printf("Completed request for %s in %d milliseconds", req.URL.Path, diff.Milliseconds())
		})
	}
}

func applyMiddlewares(logger logPrinter, mux muxHandler, requestsServed incrementer) {
	counterMiddleware := counterMiddlewareBuilder(requestsServed)
	loggingMiddleware := loggingMiddlewareBuilder(logger)

	mux.Handle("/hello", loggingMiddleware(counterMiddleware(helloHandlerBuilder(logger))))
	mux.Handle("/metrics", promhttp.Handler())
}
func SetupServer(logger logPrinter, mux muxHandler, requestsServed incrementer) {
	applyMiddlewares(logger, mux, requestsServed)
}
