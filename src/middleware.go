package main

import (
	"net/http"

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
