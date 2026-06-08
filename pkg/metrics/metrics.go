// Package metrics provides a tiny, dependency-free request metrics collector
// and a Prometheus text-exposition handler. It uses only the Go standard
// library: counters are guarded by a mutex and rendered on demand.
package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Collector records HTTP request metrics and exposes them in Prometheus text
// exposition format. The zero value is not ready for use; call NewCollector.
type Collector struct {
	mu            sync.Mutex
	total         int64
	byStatus      map[int]int64
	inFlight      int64
	durationSum   float64
	durationCount int64
}

// NewCollector returns a ready-to-use Collector.
func NewCollector() *Collector {
	return &Collector{byStatus: make(map[int]int64)}
}

// statusRecorder captures the response status code for the duration metric.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rec *statusRecorder) WriteHeader(code int) {
	if !rec.wroteHeader {
		rec.status = code
		rec.wroteHeader = true
	}
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *statusRecorder) Write(b []byte) (int, error) {
	if !rec.wroteHeader {
		rec.WriteHeader(http.StatusOK)
	}
	return rec.ResponseWriter.Write(b)
}

// Instrument wraps next so each request increments the total and per-status
// counters, tracks the in-flight gauge and accumulates request duration.
func (c *Collector) Instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.mu.Lock()
		c.inFlight++
		c.mu.Unlock()

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		elapsed := time.Since(start).Seconds()

		c.mu.Lock()
		c.inFlight--
		c.total++
		c.byStatus[rec.status]++
		c.durationSum += elapsed
		c.durationCount++
		c.mu.Unlock()
	})
}

// Handler returns an http.Handler that writes the collected metrics in
// Prometheus text exposition format.
func (c *Collector) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.mu.Lock()
		total := c.total
		inFlight := c.inFlight
		durationSum := c.durationSum
		durationCount := c.durationCount
		statuses := make([]int, 0, len(c.byStatus))
		counts := make(map[int]int64, len(c.byStatus))
		for status, count := range c.byStatus {
			statuses = append(statuses, status)
			counts[status] = count
		}
		c.mu.Unlock()
		sort.Ints(statuses)

		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintln(w, "# HELP http_requests_total Total number of HTTP requests handled.")
		fmt.Fprintln(w, "# TYPE http_requests_total counter")
		fmt.Fprintf(w, "http_requests_total %d\n", total)
		for _, status := range statuses {
			fmt.Fprintf(w, "http_requests_total{status=%q} %d\n", strconv.Itoa(status), counts[status])
		}

		fmt.Fprintln(w, "# HELP http_requests_in_flight Number of HTTP requests currently being served.")
		fmt.Fprintln(w, "# TYPE http_requests_in_flight gauge")
		fmt.Fprintf(w, "http_requests_in_flight %d\n", inFlight)

		fmt.Fprintln(w, "# HELP http_request_duration_seconds HTTP request duration in seconds.")
		fmt.Fprintln(w, "# TYPE http_request_duration_seconds summary")
		fmt.Fprintf(w, "http_request_duration_seconds_sum %g\n", durationSum)
		fmt.Fprintf(w, "http_request_duration_seconds_count %d\n", durationCount)
	})
}
