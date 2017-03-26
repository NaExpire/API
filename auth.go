package main

import (
	"fmt"
	"io"
	"net/http"
)

type businessLoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func businessLoginHandler(writer http.ResponseWriter, request *http.Request) {
	x := &businessLoginCredentials{}
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
