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
		reqID := GetRequestID(r)
		startTime := time.Now()
		lrw := newLoggingResponseWriter(w)
		h.ServeHTTP(lrw, r)
		statusCode := lrw.statusCode
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		log.Printf("request %s started at %s\n", reqID, startTime)
		log.Printf("request %s url:%s method:%s\n", reqID, r.URL, r.Method)
		log.Printf("request %s duration: %s, with status code: %d\n", reqID, elapsed, statusCode)
	})
}
