package webapi_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"promexperiment/webapi"
	"strings"
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
	origCounter := incrementer.val
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	if incrementer.val != (origCounter + 1) {
		t.Errorf("Incrementer middleware failed: expected counter to be %d, got %d", (origCounter + 1), incrementer.val)
	}
	logOutputString := logger.output.String()
	if !strings.Contains(logOutputString, "Starting request for") {
		t.Errorf("Expected logging middleware to log \"Starting request for\", log string was %q", logOutputString)
	}
	if !strings.Contains(logOutputString, "Completed request for") {
		t.Errorf("Expected logging middleware to log \"Completed request for\", log string was %q", logOutputString)
	}

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
