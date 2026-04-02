// Package api provides a framework-free HTTP handler for the GOBL API.
// It wraps the internal/cli functions and exposes them as standard
// net/http endpoints, suitable for use in any HTTP server.
package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/invopop/gobl"
	goblmcp "github.com/invopop/gobl/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// versionPrefix is the major-version path prefix derived from gobl.VERSION,
// e.g. "v0" for "v0.400.0-rc1".
var versionPrefix = "v" + fmt.Sprintf("%d", gobl.VERSION.Semver().Major())

// NewHandler builds and returns an http.Handler with all GOBL API routes
// registered under the major-version prefix (e.g. /v0/build).
// The handler includes CORS and version-header middleware.
func NewHandler() http.Handler {
	p := "/" + versionPrefix

	mux := http.NewServeMux()

	// Document operations
	mux.HandleFunc("POST "+p+"/build", handleBuild)
	mux.HandleFunc("POST "+p+"/validate", handleValidate)
	mux.HandleFunc("POST "+p+"/correct", handleCorrect)
	mux.HandleFunc("POST "+p+"/replicate", handleReplicate)
	mux.HandleFunc("POST "+p+"/sign", handleSign)
	mux.HandleFunc("POST "+p+"/verify", handleVerify)

	// Reference data (static per version — ETag for caching)
	mux.HandleFunc("GET "+p+"/schemas", withETag(handleSchemaList))
	mux.HandleFunc("GET "+p+"/schemas/{path...}", withETag(handleSchema))
	mux.HandleFunc("GET "+p+"/regimes", withETag(handleRegimeList))
	mux.HandleFunc("GET "+p+"/regimes/{code}", withETag(handleRegime))
	mux.HandleFunc("GET "+p+"/addons", withETag(handleAddonList))
	mux.HandleFunc("GET "+p+"/addons/{key...}", withETag(handleAddon))

	// MCP (Streamable HTTP — stateless)
	mcpSrv := goblmcp.NewServer()
	mcpHTTP := server.NewStreamableHTTPServer(mcpSrv, server.WithStateLess(true))
	mux.Handle("POST "+p+"/mcp", mcpHTTP)
	mux.Handle("GET "+p+"/mcp", mcpHTTP)
	mux.Handle("DELETE "+p+"/mcp", mcpHTTP)

	// Key generation
	mux.HandleFunc("POST "+p+"/keygen", handleKeygen)

	// Editor UI
	mux.HandleFunc("GET /{$}", handleEditor)

	// Favicon
	mux.HandleFunc("GET /favicon.svg", handleFavicon)

	// OpenAPI spec
	mux.HandleFunc("GET "+p+"/openapi.json", handleOpenAPI)

	// Version / health
	mux.HandleFunc("GET "+p+"/", withETag(handleVersion))

	return withCORS(withVersion(withLogging(withTiming(mux))))
}

// etag is the ETag value for all static reference-data responses,
// quoted per RFC 7232.
var etag = `"` + string(gobl.VERSION) + `"`

// withETag wraps a handler func to set an ETag header based on the GOBL
// version and return 304 Not Modified when the client already has the
// current version cached.
func withETag(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "public, max-age=86400")
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		h(w, r)
	}
}

// withVersion adds a GOBL-Version header to every response.
func withVersion(next http.Handler) http.Handler {
	version := string(gobl.VERSION)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("GOBL-Version", version)
		next.ServeHTTP(w, r)
	})
}

// withTiming adds a Server-Timing header with the request processing duration.
func withTiming(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := nowMs()
		tw := &timingWriter{ResponseWriter: w, startMs: start}
		next.ServeHTTP(tw, r)
	})
}

type timingWriter struct {
	http.ResponseWriter
	startMs     float64
	wroteHeader bool
}

func (tw *timingWriter) writeTimingHeader() {
	if !tw.wroteHeader {
		tw.wroteHeader = true
		dur := nowMs() - tw.startMs
		tw.ResponseWriter.Header().Set("Server-Timing", fmt.Sprintf("total;dur=%.4f", dur))
	}
}

func (tw *timingWriter) WriteHeader(code int) {
	tw.writeTimingHeader()
	tw.ResponseWriter.WriteHeader(code)
}

func (tw *timingWriter) Write(b []byte) (int, error) {
	tw.writeTimingHeader()
	return tw.ResponseWriter.Write(b)
}

// withLogging logs each request with method, path, status code, and duration.
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := nowMs()
		lw := &loggingWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		dur := nowMs() - start
		log.Printf("%s %s %d %.1fms", r.Method, r.URL.RequestURI(), lw.status, dur)
	})
}

type loggingWriter struct {
	http.ResponseWriter
	status int
}

func (lw *loggingWriter) WriteHeader(code int) {
	lw.status = code
	lw.ResponseWriter.WriteHeader(code)
}

// withCORS wraps a handler with permissive CORS headers.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
