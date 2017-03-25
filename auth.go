package main

import (
	"fmt"
	"io"
	"net/http"
)

type loginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		io.WriteString(writer, x.Email+"\n")
		io.WriteString(writer, x.Password+"\n")
	}
}
