package ops

import (
	"log/slog"
	"net/http"
	"time"
)

// statusRecorder is an http.ResponseWriter that remembers the first
// status code written so an outer middleware can emit it on the access
// log entry.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader records the status code and forwards to the underlying
// writer. The first call wins; subsequent calls are tracked by the
// wrapped writer but our status field reflects the initial response.
func (r *statusRecorder) WriteHeader(code int) {
	if r.status == 0 {
		r.status = code
	}
	r.ResponseWriter.WriteHeader(code)
}

// Write implements http.ResponseWriter. When the handler writes the
// body without first calling WriteHeader, net/http implicitly responds
// with 200; record that so the access log reflects it.
func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

// accessLog wraps next so every request emits one structured
// "http_request" entry on log after the handler returns. Handler-
// specific entries (e.g. who.rejected) are emitted from inside the
// handlers themselves; this baseline guarantees a record of every
// request regardless of the handler path taken.
func accessLog(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		log.Info("http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"host", stripPort(r.Host),
			"remote", r.RemoteAddr,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
