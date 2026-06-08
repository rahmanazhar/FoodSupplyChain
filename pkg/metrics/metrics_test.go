package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInstrumentAndHandler(t *testing.T) {
	c := NewCollector()

	// Drive an instrumented handler that returns a couple of statuses.
	ok := c.Instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	notFound := c.Instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	for i := 0; i < 2; i++ {
		ok.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	}
	notFound.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/missing", nil))

	rec := httptest.NewRecorder()
	c.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	body := rec.Body.String()
	wantSubstrings := []string{
		"http_requests_total 3",
		`http_requests_total{status="200"} 2`,
		`http_requests_total{status="404"} 1`,
		"http_requests_in_flight 0",
		"http_request_duration_seconds_count 3",
		"http_request_duration_seconds_sum",
	}
	for _, want := range wantSubstrings {
		if !strings.Contains(body, want) {
			t.Errorf("metrics output missing %q\n---\n%s", want, body)
		}
	}
}
