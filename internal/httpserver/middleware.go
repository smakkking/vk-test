package httpserver

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ContextKey string

// ContextKeyRequestID is the ContextKey for RequestID
const ContextKeyRequestID ContextKey = "requestID"

func reqIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()

		ctx = context.WithValue(ctx, ContextKeyRequestID, id.String())

		r = r.WithContext(ctx)

		logrus.Infof("Incomming request on %s with requestID=%s", r.RequestURI, id.String())

		next(w, r)
		// что-то делаю потом
	}
}
