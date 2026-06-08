package httpx

import (
	"net/http"
	"sync"
	"time"
)

// tokenBucket is a single client's token bucket. Tokens refill continuously at
// rps tokens per second up to burst, and are deducted one per request.
type tokenBucket struct {
	tokens   float64
	lastSeen time.Time
}

// rateLimiter is a per-client-IP token-bucket limiter guarded by a mutex. It is
// hand-rolled (no external dependencies) and refills lazily on each request.
type rateLimiter struct {
	rps   float64
	burst float64
	now   func() time.Time

	mu      sync.Mutex
	buckets map[string]*tokenBucket
}

// newRateLimiter builds a limiter allowing rps requests per second per client
// with a burst capacity. Non-positive values are clamped to sane minimums.
func newRateLimiter(rps float64, burst int) *rateLimiter {
	if rps <= 0 {
		rps = 1
	}
	if burst <= 0 {
		burst = 1
	}
	return &rateLimiter{
		rps:     rps,
		burst:   float64(burst),
		now:     time.Now,
		buckets: make(map[string]*tokenBucket),
	}
}

// allow reports whether a request from key may proceed, deducting a token when
// it does.
func (l *rateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.buckets[key]
	if !ok {
		// New clients start with a full bucket, minus the current request.
		l.buckets[key] = &tokenBucket{tokens: l.burst - 1, lastSeen: now}
		return true
	}

	// Refill based on elapsed time since the bucket was last touched.
	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens += elapsed * l.rps
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// RateLimit returns middleware enforcing a per-client-IP token-bucket limit.
// On exceeding the limit it responds 429 with a JSON error. Apply it to
// specific routes (e.g. login/register) rather than globally.
func RateLimit(rps float64, burst int) func(http.Handler) http.Handler {
	limiter := newRateLimiter(rps, burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.allow(clientIP(r)) {
				writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
