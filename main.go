package main

import (
	"encoding/json"
	"io"
	"net/http"

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
	businessRouter.HandleFunc("/login/", businessLoginHandler)
	businessRouter.HandleFunc("/register/", businessRegistrationHandler)

}

func decodeJSON(src io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(src)
	err := decoder.Decode(dst)
	return err
}

func encodeJSON(dst io.Writer, src interface{}) error {
	encoder := json.NewEncoder(dst)
	err := encoder.Encode(src)
	return err
}
