package webapi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"promexperiment/webapi"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	logger := &mockLogger{output: bytes.NewBuffer([]byte{})}
	handler := webapi.HelloHandlerBuilder(logger)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)
}
