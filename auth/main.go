package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/login", LoginHandler)
	http.ListenAndServe(":8000", nil)
}

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func DecodeJSON(src io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(src)
	err := decoder.Decode(dst)
	return err
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	x := &loginCredentials{}
	err := DecodeJSON(request.Body, x)
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
	} else {
		io.WriteString(writer, x.Username+"\n")
		io.WriteString(writer, x.Password+"\n")
	}
}
