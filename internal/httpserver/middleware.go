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

// middleware для проверки, аутентифицирован ли пользователь, вызывается при каждом дергании ресурса
func (h *HTTPService) authenticateAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok {
			correct, role := h.verifyUserPass(user, pass)

			if !correct {
				w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if role == "admin" {
				next(w, r)
			} else {
				logrus.Errorf("attempted illegal access %s to %s", user, r.URL.Path)
				http.Error(w, "you have no rights to do that", http.StatusForbidden)
			}
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}

func (h *HTTPService) authenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok {
			correct, _ := h.verifyUserPass(user, pass)

			if !correct {
				w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}

func (h *HTTPService) verifyUserPass(name string, password string) (bool, string) {
	for _, u := range h.rolesStore {
		if u.Name == name && u.Password == password {
			return true, u.Role
		}
	}
	return false, ""
}
