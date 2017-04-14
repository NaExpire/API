package business

import (
	"fmt"
	"io"
	"net/http"

	"github.com/NAExpire/API/src/util"
)

type menuItem struct {
}

type menuUpdate struct {
}

func RestaurantInfoGetHandler(writer http.ResponseWriter, request *http.Request) {
	x := &businessLoginCredentials{}
	err := util.DecodeJSON(request.Body, x)
	fmt.Printf("Got %s request to RestaurantInfoGetHandler\n", request.Method)
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
	} else {

	}
}

func RestaurantInfoUpdateHandler(writer http.ResponseWriter, request *http.Request) {
	x := &businessLoginCredentials{}
	err := util.DecodeJSON(request.Body, x)
	fmt.Printf("Got %s request to RestaurantInfoUpdateHandler\n", request.Method)
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
	} else {

	}
}

func MenuGetHandler(writer http.ResponseWriter, request *http.Request) {
	x := &businessLoginCredentials{}
	err := util.DecodeJSON(request.Body, x)
	fmt.Printf("Got %s request to MenuGetHandler\n", request.Method)
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
	} else {

	}
}

func MenuUpdateHandler(writer http.ResponseWriter, request *http.Request) {
	x := &businessLoginCredentials{}
	err := util.DecodeJSON(request.Body, x)
	fmt.Printf("Got %s request to MenuUpdateHandler\n", request.Method)
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	} else if err != nil {
		io.WriteString(writer, err.Error()+"\n")
	} else {
		
	}
}
