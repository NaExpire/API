package main

import "net/http"

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register_business", businessRegistrationHandler)
	http.ListenAndServe(":8000", nil)
}


func decodeJSON(src io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(src)
	err := decoder.Decode(dst)
	return err
}