package main

import (
	"net/http"

	"database/sql"

	"github.com/NAExpire/API/src/seshin"
	"github.com/NAExpire/API/src/util"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		h = middleware(h)
	}
	return h
}

func AllowCORS() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		})
	}
}

func DetectJsonContentType() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			// Fail fast if Content-Type is not JSON
			if contentType != "application/json" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				util.WriteErrorJSON(w, "Content-Type not recognized; should be 'application/json'")
				return
			}
			// Otherwise, continue the chain
			h.ServeHTTP(w, r)
		})
	}
}

func Authenticate(db *sql.DB, userType string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.Header.Get("session")
			if sessionID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				util.WriteErrorJSON(w, "No session header found")
				return
			}

			validated, err := seshin.ValidateSessionAndUserType(db, sessionID, userType)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				util.WriteErrorJSON(w, err.Error())
				return
			} else if !validated {
				w.WriteHeader(http.StatusForbidden)
				util.WriteErrorJSON(w, "Invalid session-id or permissions")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
