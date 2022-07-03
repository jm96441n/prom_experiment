package webapi_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"promexperiment/webapi"
	"testing"
)

func TestSetupServer(t *testing.T) {
	out := bytes.NewBuffer([]byte{})
	logger := &mockLogger{output: out}
	incrementer := &mockIncrementer{}
	m := http.NewServeMux()
	webapi.SetupServer(logger, m, incrementer)
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
}

type mockLogger struct {
	output *bytes.Buffer
}

func (m *mockLogger) Print(in ...any) {
	s := ""
	for _, a := range in {
		partial_in, ok := a.(string)
		if !ok {
			panic("expected string input for test")
		}
		s += partial_in
	}

	m.output.Write([]byte(s))
}

func (m *mockLogger) Printf(in string, args ...any) {
	s := fmt.Sprintf(in, args...)
	m.output.Write([]byte(s))
}

type mockIncrementer struct {
	val int
}

func (m *mockIncrementer) Inc() {
	m.val += 1
}
