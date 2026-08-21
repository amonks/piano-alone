package server

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// gzipMiddleware compresses responses for a caller that asked for it.
// It is the standalone server's; a host with its own compression
// middleware mounts Handler and never sees this.
//
// It stands aside for two kinds of request, both of which it used to
// break. A protocol upgrade is not a response to compress — and the
// wrapper below is not a hijackable ResponseWriter, so wrapping a
// websocket handshake turns it into a 500 that says "not implemented".
// A range request is not one either: the range was computed against
// the identity encoding, and answering it with a gzip stream under a
// 206 and the original Content-Range is a corrupt response.
//
// Upstream this wrapped two page routes, so neither case could arise;
// putting it in front of the whole mux is what made them reachable.
func gzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Vary regardless of what this request accepts: the response
		// varies by the header, so a cache that saw an uncompressed
		// answer must not hand it to a client that would have got a
		// compressed one, or the other way round.
		w.Header().Add("Vary", "Accept-Encoding")

		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
			isUpgrade(req) || req.Header.Get("Range") != "" {
			h.ServeHTTP(w, req)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		// A compressed body is a different length than the handler
		// thinks it is writing.
		w.Header().Del("Content-Length")

		gz := gzip.NewWriter(w)
		defer gz.Close()
		h.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, gz: gz}, req)
	})
}

func isUpgrade(req *http.Request) bool {
	for _, v := range req.Header.Values("Connection") {
		if strings.Contains(strings.ToLower(v), "upgrade") {
			return true
		}
	}
	return false
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) { return w.gz.Write(b) }

// Unwrap lets http.NewResponseController reach the real writer, which
// is how a handler underneath this gets at Flush, Hijack, and the
// deadline setters. Without it the standard library's controller
// stops here and reports "not implemented".
func (w *gzipResponseWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
