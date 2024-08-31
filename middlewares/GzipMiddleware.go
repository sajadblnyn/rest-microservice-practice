package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type Gzip struct {
	nextHandler http.Handler
}

func GzipMiddlewareFactory(next http.Handler) http.Handler {
	return &Gzip{nextHandler: next}
}

func (g *Gzip) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		gWriter := NewGzipResponseWriter(rw)
		gWriter.Header().Set("Content-Encoding", "gzip")
		defer gWriter.Flush()
		g.nextHandler.ServeHTTP(gWriter, r)

		return

	}

	g.nextHandler.ServeHTTP(rw, r)

}

type GzipResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewGzipResponseWriter(rw http.ResponseWriter) *GzipResponseWriter {
	return &GzipResponseWriter{rw: rw, gw: gzip.NewWriter(rw)}
}

func (gr *GzipResponseWriter) WriteHeader(statusCode int) {
	gr.rw.WriteHeader(statusCode)
}
func (gr *GzipResponseWriter) Write(d []byte) (int, error) {
	return gr.gw.Write(d)
}
func (gr *GzipResponseWriter) Header() http.Header {
	return gr.rw.Header()
}

func (gr *GzipResponseWriter) Flush() {
	gr.gw.Flush()
	gr.gw.Close()
}
