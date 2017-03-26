package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/business/login", loginHandler)
	http.HandleFunc("/business/register", businessRegistrationHandler)
	http.HandleFunc("/business/")
	http.ListenAndServe(":8000", nil)
}

func decodeJSON(src io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(src)
	err := decoder.Decode(dst)
	return err
}
