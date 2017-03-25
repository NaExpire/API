package main

import (
	"fmt"
	"io"
	"net/http"
)

type businessRegistrationCredentials struct {
	RestaurantName string `json:"restaurantName"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phoneNumber"`
	Description    string `json:"description"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

func businessRegistrationHandler(writer http.ResponseWriter, request *http.Request) {
	x := &businessRegistrationCredentials{}
	err := decodeJSON(request.Body, x)
	fmt.Printf("Got %s request to BusinessRegistrationHandler\n", request.Method)
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
	} else {
		io.WriteString(writer, x.Username+"\n")
		io.WriteString(writer, x.Password+"\n")
	}
}
