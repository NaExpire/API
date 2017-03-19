package main

import (
	"io"
	"net/http"
)

func main() {
	http.Handle("/login", &LoginHandler{})
	http.ListenAndServe(":8000", nil)
}

type LoginHandler struct {
}

func (h *LoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		io.WriteString(writer, "Goodbye world")
	}
}
