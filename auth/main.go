package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8000", nil)
}

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func decodeJSON(src io.Reader, dst interface{}) error {
	decoder := json.NewDecoder(src)
	err := decoder.Decode(dst)
	return err
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	x := &loginCredentials{}
	err := decodeJSON(request.Body, x)
	fmt.Printf("Got %s request to LoginHandler\n", request.Method)
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
	} else {
		io.WriteString(writer, x.Username+"\n")
		io.WriteString(writer, x.Password+"\n")
	}
}
