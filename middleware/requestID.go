package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const ContextKeyRequestID string = "requestID"

func RequestIDMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyRequestID, id.String())
		r = r.WithContext(ctx)

		log.Printf("Generated id %s for incomming request %s", id.String(), r.URL)
		w.Header().Set(ContextKeyRequestID, id.String())
		h.ServeHTTP(w, r)
	})
}

func GetRequestID(r *http.Request) string {
	ctx := r.Context()
	reqIDRaw := ctx.Value(ContextKeyRequestID)
	reqID, ok := reqIDRaw.(string)
	if !ok {
		log.Printf("Failed to fetch request id\n")
		return ""
	}

	return reqID
}
