package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDGeneratesAndEchoes(t *testing.T) {
	var seen string
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = RequestIDFrom(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if seen == "" {
		t.Fatal("request ID was not stored in the context")
	}
	if got := rec.Header().Get(requestIDHeader); got != seen {
		t.Fatalf("response header %q = %q, want %q", requestIDHeader, got, seen)
	}
}

func TestRequestIDHonoursInboundHeader(t *testing.T) {
	const want = "abc-123"
	var seen string
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = RequestIDFrom(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(requestIDHeader, want)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if seen != want {
		t.Fatalf("context request ID = %q, want %q", seen, want)
	}
	if got := rec.Header().Get(requestIDHeader); got != want {
		t.Fatalf("response header = %q, want %q", got, want)
	}
}

func TestRateLimitAllowsBurstThenBlocks(t *testing.T) {
	// rps is tiny so refill is negligible during the test; burst of 3 means the
	// first three requests pass and the fourth is rejected.
	h := RateLimit(0.0001, 3)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	doReq := func() int {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}

	for i := 0; i < 3; i++ {
		if got := doReq(); got != http.StatusOK {
			t.Fatalf("request %d status = %d, want 200", i+1, got)
		}
	}
	if got := doReq(); got != http.StatusTooManyRequests {
		t.Fatalf("over-limit request status = %d, want 429", got)
	}
}

func TestRateLimitIsPerClientIP(t *testing.T) {
	h := RateLimit(0.0001, 1)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	doReq := func(ip string) int {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = ip + ":1234"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}

	if got := doReq("10.0.0.1"); got != http.StatusOK {
		t.Fatalf("first client status = %d, want 200", got)
	}
	// A different client still has a full bucket.
	if got := doReq("10.0.0.2"); got != http.StatusOK {
		t.Fatalf("second client status = %d, want 200", got)
	}
	// The first client is now exhausted.
	if got := doReq("10.0.0.1"); got != http.StatusTooManyRequests {
		t.Fatalf("first client repeat status = %d, want 429", got)
	}
}
