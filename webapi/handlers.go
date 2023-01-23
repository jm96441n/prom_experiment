package webapi

import (
	"io"
	"net/http"
)

type builderFn func(logPrinter) http.Handler

var EndpointBuilders = map[string]builderFn{
	"/hello": HelloHandlerBuilder,
}

func HelloHandlerBuilder(logger logPrinter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		logger.Print("Received request for hello world")
		io.WriteString(w, "Hello World\n")
	})
}
