package middleware

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		lrw := newLoggingResponseWriter(w)
		h.ServeHTTP(lrw, r)
		statusCode := lrw.statusCode
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		log.Printf("request started at %s\n", startTime)
		log.Printf("request url:%s method:%s\n", r.URL, r.Method)
		log.Printf("request duration: %s, with status code: %d\n", elapsed, statusCode)
	})
}
