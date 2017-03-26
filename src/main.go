package main

import (
	"net/http"

	. "github.com/NAExpire/API/src/business"
	"github.com/gorilla/mux"
)

func main() {
	apiRouter := mux.NewRouter().
		StrictSlash(false)

	initBusinessRouter(apiRouter)

	http.ListenAndServe(":8000", apiRouter)
}

func initBusinessRouter(parent *mux.Router) {
	businessRouter := parent.PathPrefix("/api/business").
		Subrouter()
	businessRouter.HandleFunc("/login/", BusinessLoginHandler)
	businessRouter.HandleFunc("/register/", BusinessRegistrationHandler)
}

func initClientRotuer(parent *mux.Router) {
	// clientRouter := parent.PathPrefix("/api/client").
	// 	Subrouter()
}
