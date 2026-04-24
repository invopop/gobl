// Package api provides a framework-free HTTP handler for the GOBL API.
// It wraps the internal/ops functions and exposes them as standard
// net/http endpoints, suitable for use in any HTTP server.
//
// Use [NewHandler] to build an [http.Handler] with all GOBL API routes.
// Options allow enabling the MCP endpoint, the built-in favicon, or
// registering custom routes on the underlying ServeMux.
package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/invopop/gobl"
	goblmcp "github.com/invopop/gobl/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// VersionPrefix is the major-version path prefix derived from gobl.VERSION,
// e.g. "v0" for "v0.400.0-rc1".
var VersionPrefix = "v" + fmt.Sprintf("%d", gobl.VERSION.Semver().Major())

type config struct {
	withMCP     bool
	withFavicon bool
	extraRoutes func(mux *http.ServeMux, prefix string)
}

// Option configures the handler returned by [NewHandler].
type Option func(*config)

// WithMCP enables the MCP (Model Context Protocol) endpoint at
// /<version>/mcp.
func WithMCP() Option {
	return func(c *config) { c.withMCP = true }
}

// WithFavicon enables the built-in /favicon.svg handler.
func WithFavicon() Option {
	return func(c *config) { c.withFavicon = true }
}

// WithRoutes allows the caller to register additional routes on the
// handler's ServeMux. The prefix argument is the version path prefix
// (e.g. "/v0") so callers can mount routes alongside the API.
func WithRoutes(fn func(mux *http.ServeMux, prefix string)) Option {
	return func(c *config) { c.extraRoutes = fn }
}

// NewHandler builds and returns an [http.Handler] with GOBL API routes
// registered under the major-version prefix (e.g. /v0/build).
//
// By default the handler includes document operations (build, validate,
// correct, replicate, sign, verify), reference data (schemas, regimes,
// addons), key generation, the OpenAPI spec, and a version/health endpoint.
// CORS, version-header, timing, and logging middleware are always applied.
//
// Use [WithMCP], [WithFavicon], and [WithRoutes] to enable additional
// functionality.
func NewHandler(opts ...Option) http.Handler {
	var cfg config
	for _, o := range opts {
		o(&cfg)
	}

	p := "/" + VersionPrefix

	mux := http.NewServeMux()

	// Document operations
	mux.HandleFunc("POST "+p+"/build", handleBuild)
	mux.HandleFunc("POST "+p+"/validate", handleValidate)
	mux.HandleFunc("POST "+p+"/correct", handleCorrect)
	mux.HandleFunc("POST "+p+"/replicate", handleReplicate)
	mux.HandleFunc("POST "+p+"/sign", handleSign)
	mux.HandleFunc("POST "+p+"/verify", handleVerify)

	// Reference data (static per version — ETag for caching)
	mux.HandleFunc("GET "+p+"/schemas", WithETag(handleSchemaList))
	mux.HandleFunc("GET "+p+"/schemas/{path...}", WithETag(handleSchema))
	mux.HandleFunc("GET "+p+"/regimes", WithETag(handleRegimeList))
	mux.HandleFunc("GET "+p+"/regimes/{code}", WithETag(handleRegime))
	mux.HandleFunc("GET "+p+"/addons", WithETag(handleAddonList))
	mux.HandleFunc("GET "+p+"/addons/{key...}", WithETag(handleAddon))

	// MCP (Streamable HTTP — stateless)
	if cfg.withMCP {
		mcpSrv := goblmcp.NewServer()
		mcpHTTP := server.NewStreamableHTTPServer(mcpSrv, server.WithStateLess(true))
		mux.Handle("POST "+p+"/mcp", mcpHTTP)
		mux.Handle("GET "+p+"/mcp", mcpHTTP)
		mux.Handle("DELETE "+p+"/mcp", mcpHTTP)
	}

	// Key generation
	mux.HandleFunc("POST "+p+"/keygen", handleKeygen)

	// Favicon
	if cfg.withFavicon {
		mux.HandleFunc("GET /favicon.svg", handleFavicon)
	}

	// Custom routes
	if cfg.extraRoutes != nil {
		cfg.extraRoutes(mux, p)
	}

	// OpenAPI spec
	mux.HandleFunc("GET "+p+"/openapi.json", handleOpenAPI)

	// Version / health
	mux.HandleFunc("GET "+p+"/", WithETag(handleVersion))

	return withCORS(withVersion(withLogging(withTiming(mux))))
}

// etag is the ETag value for all static reference-data responses,
// quoted per RFC 7232.
var etag = `"` + string(gobl.VERSION) + `"`

// WithETag wraps a handler func to set an ETag header based on the GOBL
// version and return 304 Not Modified when the client already has the
// current version cached.
func WithETag(h http.HandlerFunc) http.HandlerFunc {
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
