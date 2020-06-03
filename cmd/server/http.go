package server

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CorsMiddleware(header http.Header) {
	header.Set("Access-Control-Allow-Methods", "GET")
	header.Set("Access-Control-Allow-Credentials", "true")
	header.Set("Access-Control-Allow-origin", "*")
}

func NoCacheCorsMiddleware(header http.Header) {
	header.Set("Cache-Control", "private, no-cache, max-age=0, must-revalidate, no-store, proxy-revalidate, s-maxage=0")
	header.Set("Pragma", "no-cache")
	header.Set("Expires", "0")
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func handlerGzipFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gw := gzip.NewWriter(w)
		defer gw.Close()

		gz := gzipResponseWriter{Writer: gw, ResponseWriter: w}
		fn(gz, r)
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
	CorsMiddleware(w.Header())
	NoCacheCorsMiddleware(w.Header())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	data, _ := r.Json()
	w.Write(data)
}

func httpServer(host string, port int) {
	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), handlerGzipFunc(handler))
}
